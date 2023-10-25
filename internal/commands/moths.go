package commands

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/kabili207/asamigo/internal/services/insects"
	"github.com/zekrotja/ken"
)

var mimeMaps = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

type Moths struct {
}

var (
	_ ken.SlashCommand = (*Moths)(nil)
)

func (c *Moths) Name() string {
	return "moth"
}

func (c *Moths) Description() string {
	return "Post a moth!"
}

func (c *Moths) Version() string {
	return "1.0.0"
}

func (c *Moths) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Moths) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{

			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "key",
			Description: "The ID key for a specific moth species",
			Required:    false,
		},
	}
}

func (c *Moths) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	var moth *insects.Insect
	cmd, ok := ctx.Options().GetByNameOptional("key")
	if ok {
		key := cmd.StringValue()
		moth = insects.GetInsect(insects.Moth, key)
		if moth == nil {
			return ctx.FollowUpError(fmt.Sprintf("Unable to find moth with the ID %v", key), "Cannot find moth").Send().Error
		}
	} else {
		moth = insects.GetRandomInsect(insects.Moth)
	}

	imageResp, err := http.Get(moth.Pictures[0].Url)
	if err != nil {
		return ctx.FollowUpError(fmt.Sprintf("Unhandled error:\r\n%v", err.Error()), "Error fetching image").Send().Error
	}

	contentType := imageResp.Header["Content-Type"][0]
	ext, ok := mimeMaps[contentType]
	if !ok {
		return ctx.FollowUpError(fmt.Sprintf("Unknown mime-type %v", contentType), "Error uploading image").Send().Error
	}

	return ctx.FollowUp(true, &discordgo.WebhookParams{
		Content: "## " + moth.Name + " *(" + strings.ToLower(moth.ScientificName) + ")*\n\n" +
			moth.Description,
		Files: []*discordgo.File{
			{
				Name:        strings.ReplaceAll(moth.ID, "/", "_") + ext,
				ContentType: contentType,
				Reader:      imageResp.Body,
			},
		},
	}).Send().Error
}
