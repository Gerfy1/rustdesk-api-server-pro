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
	
	// Check if user has a Personal AddressBook in the new system
	var personalAb model.AddressBook
	hasPersonalAb, err := c.Db.Where("user_id = ? AND name = ?", user.Id, model.PersonalAddressBookName).Get(&personalAb)
	if err != nil {
		return mvc.Response{
			Object: iris.Map{
				"error": err.Error(),
			},
		}
	}

	// If Personal AddressBook exists, use it; otherwise fall back to old Tags system
	var tags []string
	var tagColors map[string]int64
	var peers []iris.Map

	if hasPersonalAb {
		// Use new AddressBook system
		tagList := make([]model.AddressBookTag, 0)
		err = c.Db.Where("ab_id = ?", personalAb.Id).Find(&tagList)
		if err != nil {
			return mvc.Response{
				Object: iris.Map{
					"error": err.Error(),
				},
			}
		}
		
		tags = make([]string, 0)
		tagColors = make(map[string]int64)
		for _, tag := range tagList {
			tags = append(tags, tag.Name)
			tagColors[tag.Name] = tag.Color
		}

		// Get peers from AddressBook
		peerList := make([]model.Peer, 0)
		err = c.Db.Where("ab_id = ?", personalAb.Id).Find(&peerList)
		if err != nil {
			return mvc.Response{
				Object: iris.Map{
					"error": err.Error(),
				},
			}
		}
		
		peers = make([]iris.Map, 0)
		for _, peer := range peerList {
			var peerTags []string
			err := json.Unmarshal([]byte(peer.Tags), &peerTags)
			if err != nil {
				peerTags = []string{}
			}
			peers = append(peers, iris.Map{
				"id":       peer.RustdeskId,
				"hash":     peer.Hash,
				"username": peer.Username,
				"hostname": peer.Hostname,
				"platform": peer.Platform,
				"alias":    peer.Alias,
				"tags":     peerTags,
			})
		}
	} else {
		// Fall back to old Tags/Peer system for backward compatibility
		tagList := make([]model.Tags, 0)
		err = c.Db.Where("user_id = ?", user.Id).Find(&tagList)
		if err != nil {
			return mvc.Response{
				Object: iris.Map{
					"error": err.Error(),
				},
			}
		}
		
		tags = make([]string, 0)
		tagColors = make(map[string]int64)
		for _, tag := range tagList {
			tags = append(tags, tag.Tag)
			colorCode, err := strconv.ParseInt(tag.Color, 10, 64)
			if err != nil {
				continue
			}
			tagColors[tag.Tag] = colorCode
		}

		peerList := make([]model.Peer, 0)
		err = c.Db.Where("user_id = ?", user.Id).Find(&peerList)
		if err != nil {
			return mvc.Response{
				Object: iris.Map{
					"error": err.Error(),
				},
			}
		}
		
		peers = make([]iris.Map, 0)
		for _, peer := range peerList {
			var peerTags []string
			err := json.Unmarshal([]byte(peer.Tags), &peerTags)
			if err != nil {
				peerTags = []string{}
			}
			peers = append(peers, iris.Map{
				"id":       peer.RustdeskId,
				"hash":     peer.Hash,
				"username": peer.Username,
				"hostname": peer.Hostname,
				"platform": peer.Platform,
				"alias":    peer.Alias,
				"tags":     peerTags,
			})
		}
	}

	tagColorsJson, err := json.Marshal(tagColors)
	if err != nil {
		return mvc.Response{
			Object: iris.Map{
				"error": err.Error(),
			},
		}
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
		q := c.Db.Table(&model.AddressBook{}).Where("shared = 1")
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
	user := c.GetUser()
	abName := c.Ctx.URLParam("name")
	
	if abName == "" {
		return mvc.Response{
			Object: iris.Map{
				"error": "Address book name is required",
			},
		}
	}

	// Find the address book by name and user
	var ab model.AddressBook
	hasAb, err := c.Db.Where("user_id = ? AND name = ?", user.Id, abName).Get(&ab)
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
		}

		peers = append(peers, iris.Map{
			"id":       peer.RustdeskId,
			"hash":     peer.Hash,
			"username": peer.Username,
			"hostname": peer.Hostname,
			"platform": peer.Platform,
			"alias":    peer.Alias,
			"tags":     peerTags,
		})
	}

	tagColorsJson, err := json.Marshal(tagColors)
	if err != nil {
		return mvc.Response{
			Object: iris.Map{
				"error": err.Error(),
			},
		}
	}

	data := iris.Map{
		"tags":       tags,
		"peers":      peers,
		"tag_colors": string(tagColorsJson),
	}

	dataStr, err := json.Marshal(data)
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

	// Convert to array
	tags := make([]iris.Map, 0)
	for name, color := range tagMap {
		tags = append(tags, iris.Map{
			"name":  name,
			"color": color,
		})
	}

	return mvc.Response{
		Object: iris.Map{
			"data": tags,
		},
	}
}
