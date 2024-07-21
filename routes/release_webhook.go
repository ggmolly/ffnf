package routes

import (
	"crypto/hmac"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ggmolly/ffnf/orm"
	"github.com/ggmolly/ffnf/types"
	"github.com/gofiber/fiber/v2"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

var okActions = map[string]interface{}{
	"created":     nil,
	"published":   nil,
	"edited":      nil,
	"prereleased": nil,
	"released":    nil,
}

func verifySecret(secret string) bool {
	hash := sha512.Sum512_256([]byte(secret))
	return hex.EncodeToString(hash[:]) == os.Getenv("WEBHOOK_SECRET")
}

func headingLevel(level int) int {
	switch level {
	case 1:
		return 72
	case 2:
		return 64
	case 3:
		return 56
	case 4:
		return 48
	case 5:
		return 40
	case 6:
		return 32
	default:
		return 24
	}
}

func transformBody(body string) string {
	var output strings.Builder
	source := []byte(body)
	reader := text.NewReader(source)
	mdAst := goldmark.DefaultParser().Parse(reader)
	mdAst.Dump(source, 1)
	var nestedListLevel int

	ast.Walk(mdAst, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering && node.Kind() != ast.KindList {
			return ast.WalkContinue, nil
		}

		switch node.Kind() {
		case ast.KindHeading:
			if node.(*ast.Heading).HasBlankPreviousLines() {
				output.WriteString("\n")
			}
			output.WriteString("<size=" + strconv.Itoa(headingLevel(node.(*ast.Heading).Level)) + ">" + string(node.Text(source)) + "</size>\n")
			return ast.WalkSkipChildren, nil
		case ast.KindText:
			output.WriteString(string(node.Text(source)))
			if node.(*ast.Text).SoftLineBreak() {
				output.WriteString("\n")
			}
			// if parent is a list item, add a new line
			if node.Parent().Parent().Kind() == ast.KindListItem {
				log.Println(string(node.Text(source)))
				output.WriteString("\n")
			}
		case ast.KindParagraph:
			output.WriteString("\n")
		case ast.KindEmphasis:
			level := node.(*ast.Emphasis).Level
			if level == 1 { // render italics as orange text
				output.WriteString("<color=#f6b93b>" + string(node.Text(source)) + "</color>")
			} else if level == 2 { // render bold as outlined text
				output.WriteString("<material=outline>" + string(node.Text(source)) + "</material>")
			} else { // anything else as red text
				output.WriteString("<color=#e55039>" + string(node.Text(source)) + "</color>")
			}
			return ast.WalkSkipChildren, nil
		case ast.KindList:
			output.WriteString("\n")
			if entering {
				nestedListLevel++
			} else {
				nestedListLevel--
			}
		case ast.KindListItem:
			output.WriteString(strings.Repeat(" ", (nestedListLevel-1)*2) + "- ")
			return ast.WalkContinue, nil
		}

		return ast.WalkContinue, nil
	})
	// why should I fix bugs when I can do this instead?
	return strings.ReplaceAll(output.String(), "\n\n", "\n")
}

// POST /api/v1/webhook/releases/:secret
func ReleaseWebhook(c *fiber.Ctx) error {
	// Verify URL's secret
	if !verifySecret(c.Params("secret")) {
		return c.Status(fiber.StatusForbidden).SendString("invalid secret")
	}

	// Check repo signature
	sharedSecret := os.Getenv("WEBHOOK_SECRET")
	signature := c.Get("X-Hub-Signature-256")

	// HMAC-SHA256 of the shared secret and the request body
	hmacContext := hmac.New(sha512.New, []byte(sharedSecret))
	hmacContext.Write([]byte(c.Body()))

	expectedSignature := "sha256=" + hex.EncodeToString(hmacContext.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(signature), []byte(expectedSignature)) != 1 {
		return c.Status(fiber.StatusForbidden).SendString("invalid signature")
	}

	var payload types.GithubReleaseWebhook
	if err := json.Unmarshal(c.Body(), &payload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid json")
	}

	if _, ok := okActions[payload.Action]; !ok {
		return c.Status(fiber.StatusBadRequest).SendString("invalid action: " + payload.Action)
	}

	release := orm.Release{
		ID:         payload.Release.ID,
		Prerelease: payload.Release.Prerelease,
		Name:       payload.Release.Name,
		Body:       payload.Release.Body,
	}

	notice := orm.Notice{
		ID:          payload.Release.ID,
		ButtonTitle: payload.Release.Name,
		Title:       payload.Release.Name,
		Content:     transformBody(payload.Release.Body),
		CreatedAt:   payload.Release.CreatedAt,
	}

	if payload.Action == "edited" {
		// Delete old release / notice to give the new one a new ID
		orm.GormDB.Delete(&release)
		orm.GormDB.Delete(&notice)
	}

	tx := orm.GormDB.Begin()
	// Create the new release
	if err := tx.Create(&release).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).SendString("failed to save release: " + err.Error())
	}
	// Create the new notice
	if err := tx.Create(&notice).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).SendString("failed to save notice: " + err.Error())
	}
	tx.Commit() // save if everything went well
	return c.SendString("ok")
}
