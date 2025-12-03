package admin

import (
	"encoding/json"
	"rustdesk-api-server-pro/app/model"
	"rustdesk-api-server-pro/config"
	"rustdesk-api-server-pro/db"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"xorm.io/xorm"
)

type AddressBooksController struct {
	basicController
}

func (c *AddressBooksController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/address-books/list", "HandleList")
	b.Handle("GET", "/address-books/{id:int}", "HandleGet")
	b.Handle("POST", "/address-books", "HandleCreate")
	b.Handle("PUT", "/address-books/{id:int}", "HandleUpdate")
	b.Handle("DELETE", "/address-books/{id:int}", "HandleDelete")
	b.Handle("GET", "/address-books/{id:int}/peers", "HandleGetPeers")
	b.Handle("GET", "/address-books/{id:int}/tags", "HandleGetTags")
	b.Handle("POST", "/address-books/{id:int}/import-devices", "HandleImportDevices")
	b.Handle("POST", "/address-books/{id:int}/peers", "HandleAddPeer")
	b.Handle("DELETE", "/address-books/{id:int}/peers/{peerId:int}", "HandleDeletePeer")
}

// HandleList - List all address books with pagination
func (c *AddressBooksController) HandleList() mvc.Result {
	currentPage := c.Ctx.URLParamIntDefault("current", 1)
	pageSize := c.Ctx.URLParamIntDefault("size", 10)
	name := c.Ctx.URLParamDefault("name", "")
	owner := c.Ctx.URLParamDefault("owner", "")

	query := func() *xorm.Session {
		q := c.Db.Table(&model.AddressBook{})

		if name != "" {
			q.Where("name LIKE ?", "%"+name+"%")
		}
		if owner != "" {
			q.Where("owner LIKE ?", "%"+owner+"%")
		}
		q.Desc("created_at")
		return q
	}

	pagination := db.NewPagination(currentPage, pageSize)
	abList := make([]model.AddressBook, 0)

	err := pagination.Paginate(query, &model.AddressBook{}, &abList)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	list := make([]iris.Map, 0)
	for _, ab := range abList {
		// Count peers in this address book
		peerCount, _ := c.Db.Where("ab_id = ?", ab.Id).Count(&model.Peer{})

		list = append(list, iris.Map{
			"id":         ab.Id,
			"user_id":    ab.UserId,
			"guid":       ab.Guid,
			"name":       ab.Name,
			"owner":      ab.Owner,
			"note":       ab.Note,
			"rule":       ab.Rule,
			"max_peer":   ab.MaxPeer,
			"shared":     ab.Shared,
			"peer_count": peerCount,
			"created_at": ab.CreatedAt.Format(config.TimeFormat),
			"updated_at": ab.UpdatedAt.Format(config.TimeFormat),
		})
	}

	return c.Success(iris.Map{
		"total":   pagination.TotalCount,
		"records": list,
		"current": currentPage,
		"size":    pageSize,
	}, "ok")
}

// HandleGet - Get single address book
func (c *AddressBooksController) HandleGet() mvc.Result {
	id, err := c.Ctx.Params().GetInt("id")
	if err != nil {
		return c.Error(nil, "Invalid ID")
	}

	var ab model.AddressBook
	has, err := c.Db.ID(id).Get(&ab)
	if err != nil {
		return c.Error(nil, err.Error())
	}
	if !has {
		return c.Error(nil, "Address book not found")
	}

	return c.Success(iris.Map{
		"id":         ab.Id,
		"user_id":    ab.UserId,
		"guid":       ab.Guid,
		"name":       ab.Name,
		"owner":      ab.Owner,
		"note":       ab.Note,
		"rule":       ab.Rule,
		"max_peer":   ab.MaxPeer,
		"shared":     ab.Shared,
		"created_at": ab.CreatedAt.Format(config.TimeFormat),
		"updated_at": ab.UpdatedAt.Format(config.TimeFormat),
	}, "ok")
}

// HandleCreate - Create new address book
func (c *AddressBooksController) HandleCreate() mvc.Result {
	var form struct {
		UserId  int      `json:"user_id"`
		Name    string   `json:"name"`
		Note    string   `json:"note"`
		Rule    int      `json:"rule"`
		MaxPeer int      `json:"max_peer"`
		Shared  bool     `json:"shared"`
		Tags    []string `json:"tags"`
	}

	if err := c.Ctx.ReadJSON(&form); err != nil {
		return c.Error(nil, "Invalid request body")
	}

	// Validate
	if form.Name == "" {
		return c.Error(nil, "Name is required")
	}

	// Get owner username
	var user model.User
	has, err := c.Db.ID(form.UserId).Get(&user)
	if err != nil || !has {
		return c.Error(nil, "User not found")
	}

	ab := model.AddressBook{
		UserId:  form.UserId,
		Name:    form.Name,
		Owner:   user.Username,
		Note:    form.Note,
		Rule:    form.Rule,
		MaxPeer: form.MaxPeer,
		Shared:  form.Shared,
	}

	_, err = c.Db.Insert(&ab)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	// Create tags if provided
	if len(form.Tags) > 0 {
		for _, tagName := range form.Tags {
			if tagName == "" {
				continue
			}
			
			tag := model.AddressBookTag{
				UserId: form.UserId,
				AbId:   ab.Id,
				Name:   tagName,
				Color:  4278190335, // Default color (blue) for RustDesk 1.4.x
			}
			inserted, err := c.Db.Insert(&tag)
			if err != nil {
				c.Ctx.Application().Logger().Errorf("Failed to insert tag %s: %v", tagName, err)
			} else {
				c.Ctx.Application().Logger().Infof("Created tag: %s (id=%d) for address book %d", tagName, tag.Id, ab.Id)
				_ = inserted
			}
		}
	}

	return c.Success(iris.Map{
		"id": ab.Id,
	}, "Address book created successfully")
}

// HandleUpdate - Update address book
func (c *AddressBooksController) HandleUpdate() mvc.Result {
	id, err := c.Ctx.Params().GetInt("id")
	if err != nil {
		return c.Error(nil, "Invalid ID")
	}

	var form struct {
		Name    string `json:"name"`
		Note    string `json:"note"`
		Rule    int    `json:"rule"`
		MaxPeer int    `json:"max_peer"`
		Shared  bool   `json:"shared"`
	}

	if err := c.Ctx.ReadJSON(&form); err != nil {
		return c.Error(nil, "Invalid request body")
	}

	var ab model.AddressBook
	has, err := c.Db.ID(id).Get(&ab)
	if err != nil || !has {
		return c.Error(nil, "Address book not found")
	}

	ab.Name = form.Name
	ab.Note = form.Note
	ab.Rule = form.Rule
	ab.MaxPeer = form.MaxPeer
	ab.Shared = form.Shared

	_, err = c.Db.ID(id).Update(&ab)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	return c.Success(nil, "Address book updated successfully")
}

// HandleDelete - Delete address book
func (c *AddressBooksController) HandleDelete() mvc.Result {
	id, err := c.Ctx.Params().GetInt("id")
	if err != nil {
		return c.Error(nil, "Invalid ID")
	}

	var ab model.AddressBook
	has, err := c.Db.ID(id).Get(&ab)
	if err != nil || !has {
		return c.Error(nil, "Address book not found")
	}

	// Delete associated peers
	c.Db.Where("ab_id = ?", id).Delete(&model.Peer{})

	// Delete address book
	_, err = c.Db.ID(id).Delete(&ab)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	return c.Success(nil, "Address book deleted successfully")
}

// HandleGetPeers - Get all peers in an address book with their online status
func (c *AddressBooksController) HandleGetPeers() mvc.Result {
	id, err := c.Ctx.Params().GetInt("id")
	if err != nil {
		return c.Error(nil, "Invalid ID")
	}

	// Get peers
	peers := make([]model.Peer, 0)
	err = c.Db.Where("ab_id = ?", id).Find(&peers)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	// Build result with device status
	list := make([]iris.Map, 0)
	for _, peer := range peers {
		// Check if device is online
		var device model.Device
		isOnline := false
		lastSeenAt := ""
		ipAddress := ""
		
		has, err := c.Db.Where("rustdesk_id = ?", peer.RustdeskId).Get(&device)
		if err == nil && has {
			isOnline = device.IsOnline
			lastSeenAt = device.LastSeenAt.Format(config.TimeFormat)
			ipAddress = device.IpAddress
		}

		list = append(list, iris.Map{
			"id":            peer.Id,
			"rustdesk_id":   peer.RustdeskId,
			"username":      peer.Username,
			"hostname":      peer.Hostname,
			"alias":         peer.Alias,
			"platform":      peer.Platform,
			"tags":          peer.Tags,
			"is_online":     isOnline,
			"last_seen_at":  lastSeenAt,
			"ip_address":    ipAddress,
			"created_at":    peer.CreatedAt.Format(config.TimeFormat),
		})
	}

	return c.Success(iris.Map{
		"total":   len(list),
		"records": list,
	}, "ok")
}

// HandleGetTags - Get all tags for an address book
func (c *AddressBooksController) HandleGetTags() mvc.Result {
	id := c.Ctx.Params().GetIntDefault("id", 0)

	// Get tags for this address book
	tagList := make([]model.AddressBookTag, 0)
	err := c.Db.Where("ab_id = ?", id).Find(&tagList)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	// Format tags as simple array of strings
	tags := make([]string, 0)
	for _, tag := range tagList {
		tags = append(tags, tag.Name)
	}

	c.Ctx.Application().Logger().Infof("HandleGetTags: Found %d tags for address book %d", len(tags), id)

	return c.Success(tags, "ok")
}

// HandleImportDevices - Import devices as peers into an address book
func (c *AddressBooksController) HandleImportDevices() mvc.Result {
	id := c.Ctx.Params().GetIntDefault("id", 0)

	// Get address book
	var ab model.AddressBook
	has, err := c.Db.Where("id = ?", id).Get(&ab)
	if err != nil {
		return c.Error(nil, err.Error())
	}
	if !has {
		return c.Error(nil, "Address book not found")
	}

	// Get online devices
	deviceList := make([]model.Device, 0)
	err = c.Db.Where("is_online = ?", true).Find(&deviceList)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	importedCount := 0
	skippedCount := 0

	for _, device := range deviceList {
		// Check if peer already exists
		var existingPeer model.Peer
		has, _ := c.Db.Where("rustdesk_id = ? AND ab_id = ?", device.RustdeskId, ab.Id).Get(&existingPeer)
		
		if has {
			skippedCount++
			continue
		}

		// Create new peer from device
		peer := model.Peer{
			UserId:     ab.UserId,
			AbId:       ab.Id,
			RustdeskId: device.RustdeskId,
			Username:   device.Username,
			Hostname:   device.Hostname,
			Platform:   device.Os,
			Alias:      device.RustdeskId, // Use RustDesk ID as alias so it shows correctly in client
			Tags:       "[]",
		}

		_, err = c.Db.Insert(&peer)
		if err != nil {
			c.Ctx.Application().Logger().Errorf("Failed to import device %s: %v", device.RustdeskId, err)
			skippedCount++
			continue
		}

		importedCount++
	}

	return c.Success(iris.Map{
		"imported": importedCount,
		"skipped":  skippedCount,
		"total":    len(deviceList),
	}, "ok")
}

// HandleAddPeer - Add a peer manually to an address book
func (c *AddressBooksController) HandleAddPeer() mvc.Result {
	abId, err := c.Ctx.Params().GetInt("id")
	if err != nil {
		return c.Error(500, "Invalid address book ID")
	}

	// Verify address book exists
	var ab model.AddressBook
	has, err := c.Db.Where("id = ?", abId).Get(&ab)
	if err != nil {
		return c.Error(500, err.Error())
	}
	if !has {
		return c.Error(404, "Address book not found")
	}

	// Parse request body
	type AddPeerForm struct {
		AbId       int    `json:"ab_id"`
		RustdeskId string `json:"rustdesk_id"`
		Alias      string `json:"alias"`
		Password   string `json:"password"`
		Hostname   string `json:"hostname"`
		Username   string `json:"username"`
		Platform   string `json:"platform"`
		Tags       string `json:"tags"`
	}

	var form AddPeerForm
	err = c.Ctx.ReadJSON(&form)
	if err != nil {
		return c.Error(400, "Invalid request body")
	}

	// Validate required fields
	if form.RustdeskId == "" {
		return c.Error(400, "RustDesk ID is required")
	}

	// Check if peer already exists
	var existingPeer model.Peer
	has, _ = c.Db.Where("rustdesk_id = ? AND ab_id = ?", form.RustdeskId, ab.Id).Get(&existingPeer)
	if has {
		return c.Error(400, "Peer with this RustDesk ID already exists in this address book")
	}

	// Create new peer
	peer := model.Peer{
		UserId:     ab.UserId,
		AbId:       ab.Id,
		RustdeskId: form.RustdeskId,
		Alias:      form.Alias,
		Password:   form.Password,
		Hostname:   form.Hostname,
		Username:   form.Username,
		Platform:   form.Platform,
		Tags:       form.Tags,
	}

	// If tags is empty, set to empty array
	if peer.Tags == "" {
		peer.Tags = "[]"
	}

	// Parse tags and create AddressBookTag entries if they don't exist
	var tagNames []string
	err = json.Unmarshal([]byte(peer.Tags), &tagNames)
	if err == nil && len(tagNames) > 0 {
		// For each tag name, check if it exists in AddressBookTag
		for _, tagName := range tagNames {
			if tagName == "" {
				continue
			}
			
			var existingTag model.AddressBookTag
			hasTag, _ := c.Db.Where("ab_id = ? AND name = ?", ab.Id, tagName).Get(&existingTag)
			
			if !hasTag {
				// Create new tag with default color
				newTag := model.AddressBookTag{
					UserId: ab.UserId,
					AbId:   ab.Id,
					Name:   tagName,
					Color:  0xFF0000FF, // Default blue color
				}
				_, _ = c.Db.Insert(&newTag)
			}
		}
	}

	_, err = c.Db.Insert(&peer)
	if err != nil {
		return c.Error(500, "Failed to create peer: "+err.Error())
	}

	return c.Success(iris.Map{
		"id":          peer.Id,
		"rustdesk_id": peer.RustdeskId,
		"alias":       peer.Alias,
	}, "Peer created successfully")
}

// HandleDeletePeer - Delete a peer from an address book
func (c *AddressBooksController) HandleDeletePeer() mvc.Result {
	abId, err := c.Ctx.Params().GetInt("id")
	if err != nil {
		return c.Error(500, "Invalid address book ID")
	}

	peerId, err := c.Ctx.Params().GetInt("peerId")
	if err != nil {
		return c.Error(500, "Invalid peer ID")
	}

	// Verify address book exists
	var ab model.AddressBook
	has, err := c.Db.Where("id = ?", abId).Get(&ab)
	if err != nil {
		return c.Error(500, err.Error())
	}
	if !has {
		return c.Error(404, "Address book not found")
	}

	// Verify peer exists and belongs to this address book
	var peer model.Peer
	has, err = c.Db.Where("id = ? AND ab_id = ?", peerId, ab.Id).Get(&peer)
	if err != nil {
		return c.Error(500, err.Error())
	}
	if !has {
		return c.Error(404, "Peer not found in this address book")
	}

	// Delete the peer
	_, err = c.Db.ID(peerId).Delete(&model.Peer{})
	if err != nil {
		return c.Error(500, "Failed to delete peer: "+err.Error())
	}

	return c.Success(iris.Map{
		"id": peerId,
	}, "Peer deleted successfully")
}

