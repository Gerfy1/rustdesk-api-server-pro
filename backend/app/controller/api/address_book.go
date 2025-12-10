package api

import (
	"encoding/json"
	"rustdesk-api-server-pro/app/form/api"
	"rustdesk-api-server-pro/app/model"
	"rustdesk-api-server-pro/db"
	"strconv"

	"github.com/beevik/guid"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"xorm.io/xorm"
)

type AddressBookController struct {
	basicController
}

func (c *AddressBookController) GetAb() mvc.Result {
	user := c.GetUser()
	
	// Get ALL tags from ALL address books (shared access)
	tagList := make([]model.AddressBookTag, 0)
	err := c.Db.Find(&tagList)
	if err != nil {
		return mvc.Response{
			Object: iris.Map{
				"error": err.Error(),
			},
		}
	}
	
	tags := make([]string, 0)
	tagColors := make(map[string]int64)
	seenTags := make(map[string]bool)
	
	for _, tag := range tagList {
		if !seenTags[tag.Name] {
			tags = append(tags, tag.Name)
			tagColors[tag.Name] = tag.Color
			seenTags[tag.Name] = true
		}
	}

	// Get ALL peers from ALL address books (shared access)
	peerList := make([]model.Peer, 0)
	err = c.Db.Find(&peerList)
	if err != nil {
		return mvc.Response{
			Object: iris.Map{
				"error": err.Error(),
			},
		}
	}
	
	peers := make([]iris.Map, 0)
	for _, peer := range peerList {
		var peerTags []string
		if peer.Tags != "" {
			err := json.Unmarshal([]byte(peer.Tags), &peerTags)
			if err != nil {
				peerTags = []string{}
			}
		} else {
			peerTags = []string{}
		}
		
		// Ensure hash is a valid string (not binary data)
		hash := peer.Hash
		if hash == "" {
			hash = ""
		}
		
		peers = append(peers, iris.Map{
			"id":       peer.RustdeskId,
			"hash":     hash,
			"username": peer.Username,
			"hostname": peer.Hostname,
			"platform": peer.Platform,
			"alias":    peer.Alias,
			"tags":     peerTags,
		})
	}
	
	// Log for debugging
	c.Ctx.Application().Logger().Infof("User %d (%s) accessing shared address books: %d tags, %d peers", 
		user.Id, user.Username, len(tags), len(peers))

	tagColorsJson, err := json.Marshal(tagColors)
	if err != nil {
		tagColorsJson = []byte("{}")
	}

	dataJson, err := json.Marshal(iris.Map{
		"tags":       tags,
		"peers":      peers,
		"tag_colors": string(tagColorsJson),
	})
	if err != nil {
		return mvc.Response{
			Object: iris.Map{
				"error": err.Error(),
			},
		}
	}

	return mvc.Response{
		Object: iris.Map{
			"licensed_devices": user.LicensedDevices,
			"data":             string(dataJson),
		},
	}
}

func (c *AddressBookController) PostAb() mvc.Result {
	var abForm api.AbForm
	err := c.Ctx.ReadJSON(&abForm)
	if err != nil {
		return mvc.Response{
			Object: iris.Map{
				"error": err.Error(),
			},
		}
	}
	var abData api.AbData
	err = json.Unmarshal([]byte(abForm.Data), &abData)
	if err != nil {
		return mvc.Response{
			Object: iris.Map{
				"error": err.Error(),
			},
		}
	}
	var tagColors map[string]int64
	err = json.Unmarshal([]byte(abData.TagColors), &tagColors)
	if err != nil {
		return mvc.Response{
			Object: iris.Map{
				"error": err.Error(),
			},
		}
	}

	session := c.Db.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return mvc.Response{
			Object: iris.Map{
				"error": err.Error(),
			},
		}
	}
	user := c.GetUser()
	if user.LicensedDevices > 0 && len(abData.Peers) > user.LicensedDevices {
		return mvc.Response{
			Object: iris.Map{
				"error": "Number of equipment in excess of licenses",
			},
		}
	}
	_, err = session.Where("user_id = ?", user.Id).Delete(&model.Tags{})
	if err != nil {
		_ = session.Rollback()
		return mvc.Response{
			Object: iris.Map{
				"error": err.Error(),
			},
		}
	}

	_, err = session.Where("user_id = ?", user.Id).Delete(&model.Peer{})
	if err != nil {
		_ = session.Rollback()
		return mvc.Response{
			Object: iris.Map{
				"error": err.Error(),
			},
		}
	}

	tags := make([]*model.Tags, 0)
	for _, tag := range abData.Tags {
		tags = append(tags, &model.Tags{
			UserId: user.Id,
			Tag:    tag,
			Color:  strconv.FormatInt(tagColors[tag], 10),
		})
	}
	if len(tags) > 0 {
		_, err = session.Insert(tags)
		if err != nil {
			_ = session.Rollback()
			return mvc.Response{
				Object: iris.Map{
					"error": err.Error(),
				},
			}
		}
	}

	peers := make([]*model.Peer, 0)
	for _, peer := range abData.Peers {
		peerTags := ""
		b, err := json.Marshal(peer.Tags)
		if err == nil {
			peerTags = string(b)
		}
		peers = append(peers, &model.Peer{
			UserId:     user.Id,
			RustdeskId: peer.Id,
			Hash:       peer.Hash,
			Username:   peer.Username,
			Hostname:   peer.Hostname,
			Platform:   peer.Platform,
			Alias:      peer.Alias,
			Tags:       peerTags,
		})
	}
	if len(peers) > 0 {
		_, err = session.Insert(peers)
		if err != nil {
			_ = session.Rollback()
			return mvc.Response{
				Object: iris.Map{
					"error": err.Error(),
				},
			}
		}
	}

	err = session.Commit()
	if err != nil {
		return mvc.Response{
			Object: iris.Map{
				"error": err.Error(),
			},
		}
	}

	return mvc.Response{}
}

func (c *AddressBookController) PostAbPersonal() mvc.Result {
	user := c.GetUser()
	var ab model.AddressBook
	has, err := c.Db.Where("user_id = ?", user.Id).Get(&ab)
	if err != nil {
		return mvc.Response{
			Object: iris.Map{
				"error": err.Error(),
			},
		}
	}

	if !has {
		g := guid.New()
		ab.UserId = user.Id
		ab.Guid = g.String()
		ab.Name = model.PersonalAddressBookName
		ab.Owner = user.Username
		ab.MaxPeer = model.MaxPeer
		ab.Note = "default address book"
		ab.Rule = 3 // full control
		c.Db.Insert(&ab)
	}

	return mvc.Response{
		Object: iris.Map{
			"guid": ab.Guid,
		},
	}
}

func (c *AddressBookController) PostAbSettings() mvc.Result {
	user := c.GetUser()
	var ab model.AddressBook
	_, _ = c.Db.Where("user_id = ?", user.Id).Get(&ab)
	return mvc.Response{
		Object: iris.Map{
			"max_peer_one_ab": ab.MaxPeer,
		},
	}
}

func (c *AddressBookController) PostAbSharedProfiles() mvc.Result {
	current := c.Ctx.URLParamIntDefault("current", 1)
	pageSize := c.Ctx.URLParamIntDefault("pageSize", 10)

	query := func() *xorm.Session {
		// Return ALL address books (shared access for all users)
		q := c.Db.Table(&model.AddressBook{})
		return q
	}

	pagination := db.NewPagination(current, pageSize)
	sharedAbList := make([]model.AddressBook, 0)
	err := pagination.Paginate(query, &model.AddressBook{}, &sharedAbList)
	if err != nil {
		return mvc.Response{
			Object: iris.Map{
				"error": err.Error(),
			},
		}
	}
	data := make([]iris.Map, 0)
	for _, ab := range sharedAbList {
		data = append(data, iris.Map{
			"guid":  ab.Guid,
			"name":  ab.Name,
			"owner": ab.Owner,
			"note":  ab.Note,
			"rule":  ab.Rule,
		})
	}

	return mvc.Response{
		Object: iris.Map{
			"total": pagination.TotalCount,
			"data":  data,
		},
	}
}

// GetAbGet handles GET /api/ab/get?name=xxx - Returns peers for a specific address book
func (c *AddressBookController) GetAbGet() mvc.Result {
	abName := c.Ctx.URLParam("name")
	
	if abName == "" {
		return mvc.Response{
			Object: iris.Map{
				"error": "Address book name is required",
			},
		}
	}

	// Find the address book by name (allow access to any address book)
	var ab model.AddressBook
	hasAb, err := c.Db.Where("name = ?", abName).Get(&ab)
	if err != nil {
		return mvc.Response{
			Object: iris.Map{
				"error": err.Error(),
			},
		}
	}

	if !hasAb {
		return mvc.Response{
			Object: iris.Map{
				"error": "Address book not found",
			},
		}
	}

	// Get tags for this address book
	tagList := make([]model.AddressBookTag, 0)
	err = c.Db.Where("ab_id = ?", ab.Id).Find(&tagList)
	if err != nil {
		return mvc.Response{
			Object: iris.Map{
				"error": err.Error(),
			},
		}
	}

	tags := make([]string, 0)
	tagColors := make(map[string]int64)
	for _, tag := range tagList {
		tags = append(tags, tag.Name)
		tagColors[tag.Name] = tag.Color
	}

	// Get peers for this address book
	peerList := make([]model.Peer, 0)
	err = c.Db.Where("ab_id = ?", ab.Id).Find(&peerList)
	if err != nil {
		return mvc.Response{
			Object: iris.Map{
				"error": err.Error(),
			},
		}
	}

	peers := make([]iris.Map, 0)
	for _, peer := range peerList {
		var peerTags []string
		if peer.Tags != "" {
			err := json.Unmarshal([]byte(peer.Tags), &peerTags)
			if err != nil {
				peerTags = []string{}
			}
		} else {
			peerTags = []string{}
		}

		// Ensure hash is a valid string (not binary data)
		hash := peer.Hash
		if hash == "" {
			hash = ""
		}

		peers = append(peers, iris.Map{
			"id":       peer.RustdeskId,
			"hash":     hash,
			"username": peer.Username,
			"hostname": peer.Hostname,
			"platform": peer.Platform,
			"alias":    peer.Alias,
			"tags":     peerTags,
		})
	}

	// Serialize tag_colors as JSON string
	tagColorsJson, err := json.Marshal(tagColors)
	if err != nil {
		tagColorsJson = []byte("{}")
	}

	// Build the data object with proper structure
	dataObj := iris.Map{
		"tags":       tags,
		"peers":      peers,
		"tag_colors": string(tagColorsJson),
	}

	dataStr, err := json.Marshal(dataObj)
	if err != nil {
		return mvc.Response{
			Object: iris.Map{
				"error": err.Error(),
			},
		}
	}

	return mvc.Response{
		Object: iris.Map{
			"data": string(dataStr),
		},
	}
}

// PostAbTags handles POST /api/ab/tags - Returns tags for address books
func (c *AddressBookController) PostAbTags() mvc.Result {
	user := c.GetUser()
	
	// Get all address books for this user
	abList := make([]model.AddressBook, 0)
	err := c.Db.Where("user_id = ?", user.Id).Find(&abList)
	if err != nil {
		return mvc.Response{
			Object: iris.Map{
				"error": err.Error(),
			},
		}
	}

	// Collect all ab_ids
	abIds := make([]int, 0)
	for _, ab := range abList {
		abIds = append(abIds, ab.Id)
	}

	// Get all tags for all address books of this user
	tagList := make([]model.AddressBookTag, 0)
	
	if len(abIds) > 0 {
		err = c.Db.In("ab_id", abIds).Find(&tagList)
		if err != nil {
			return mvc.Response{
				Object: iris.Map{
					"error": err.Error(),
				},
			}
		}
	}

	// Format tags as array of objects with name and color
	// Use map to avoid duplicates
	tagMap := make(map[string]int64)
	for _, tag := range tagList {
		tagMap[tag.Name] = tag.Color
	}

	// Convert to array of structs (proper JSON structure for RustDesk 1.4.x)
	type TagResponse struct {
		Name  string `json:"name"`
		Color int64  `json:"color"`
	}
	
	tags := make([]TagResponse, 0)
	for name, color := range tagMap {
		tags = append(tags, TagResponse{
			Name:  name,
			Color: color,
		})
	}

	// Log for debugging
	c.Ctx.Application().Logger().Infof("PostAbTags: Found %d tags for user %d", len(tags), user.Id)

	// Return array directly without wrapping in "data" key
	// RustDesk 1.4.x expects: [{"name":"tag1","color":4278190335}, ...]
	return mvc.Response{
		Object: tags,
	}
}
