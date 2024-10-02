package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/internal/blocks"
	"github.com/nathanhollows/Rapua/internal/flash"
	templates "github.com/nathanhollows/Rapua/internal/templates/admin"
	bTemplates "github.com/nathanhollows/Rapua/internal/templates/blocks"
)

// BlockEdit shows the form to edit a block
func (h *AdminHandler) BlockEdit(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())
	location := chi.URLParam(r, "location")
	if !h.GameManagerService.ValidateLocationID(user, location) {
		h.Logger.Error("BlockEdit: invalid location", "location", location)
		err := templates.Toast(*flash.NewError("Could not find block")).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("BlockEdit: rendering template", "error", err)
		}
		return
	}

	blockID := chi.URLParam(r, "blockID")

	block, err := h.BlockService.GetByBlockID(r.Context(), blockID)
	if err != nil {
		h.Logger.Error("BlockEdit: getting block", "error", err)
		err := templates.Toast(*flash.NewError("Could not find block")).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("BlockEdit: rendering template", "error", err)
		}
		return
	}

	if block.GetLocationID() != location {
		h.Logger.Error("BlockEdit: block does not belong to location", "blockID", blockID, "location", location)
		err := templates.Toast(*flash.NewError("Could not find block")).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("BlockEdit: rendering template", "error", err)
		}
		return
	}

	w.Header().Set("Content-Type", "text/html")
	switch block.(type) {
	case *blocks.MarkdownBlock:
		b := block.(*blocks.MarkdownBlock)
		err = bTemplates.MarkdownAdmin(user.CurrentInstance.Settings, *b).Render(r.Context(), w)
	case *blocks.PasswordBlock:
		b := block.(*blocks.PasswordBlock)
		err = bTemplates.PasswordAdmin(user.CurrentInstance.Settings, *b).Render(r.Context(), w)
	}

}

// BlockEditPost updates the block
func (h *AdminHandler) BlockEditPost(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	location := chi.URLParam(r, "location")
	if !h.GameManagerService.ValidateLocationID(user, location) {
		h.Logger.Error("BlockEditPost: invalid location", "location", location)
		err := templates.Toast(*flash.NewError("Could not update block. Invalid location")).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("BlockEditPost: rendering template", "error", err)
		}
		return
	}

	blockID := chi.URLParam(r, "blockID")

	block, err := h.BlockService.GetByBlockID(r.Context(), blockID)
	if err != nil {
		h.Logger.Error("BlockEditPost: getting block", "error", err)
		err := templates.Toast(*flash.NewError("Could not update block")).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("BlockEditPost: rendering template", "error", err)
		}
		return
	}

	if block.GetLocationID() != location {
		h.Logger.Error("BlockEditPost: block does not belong to location", "blockID", blockID, "location", location)
		err := templates.Toast(*flash.NewError("Could not update block")).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("BlockEditPost: rendering template", "error", err)
		}
		return
	}

	err = h.BlockService.UpdateBlock(r.Context(), &block)
	if err != nil {
		h.Logger.Error("BlockEditPost: updating block", "error", err)
		err := templates.Toast(*flash.NewError("Could not update block")).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("BlockEditPost: rendering template", "error", err)
		}
		return
	}

	err = templates.Toast(*flash.NewSuccess("Block updated")).Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("BlockEditPost: rendering template", "error", err)
	}
}

// Show the form to edit the navigation settings.
func (h *AdminHandler) BlockNewPost(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	blockType := chi.URLParam(r, "type")

	location := chi.URLParam(r, "location")
	if !h.GameManagerService.ValidateLocationID(user, location) {
		h.Logger.Error("BlockNewPost: invalid location", "location", location)
		err := templates.Toast(*flash.NewError("Could not create block. Invalid location")).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("BlockNewPost: rendering template", "error", err)
		}
		return
	}

	block, err := h.BlockService.NewBlock(r.Context(), location, blockType)
	if err != nil {
		h.Logger.Error("BlockNewPost: creating block", "error", err)
		err := templates.Toast(*flash.NewError("Could not create block")).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("BlockNewPost: rendering template", "error", err)
		}
		return
	}

	data := block.GetAdminData()
	switch data.(type) {
	case blocks.MarkdownBlock:
		b := data.(blocks.MarkdownBlock)
		err = bTemplates.MarkdownAdmin(b).Render(r.Context(), w)
	}

}
