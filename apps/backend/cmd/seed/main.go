package main

import (
	"backend/internal/crypto"
	"backend/internal/database"
	"backend/internal/types"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	random "math/rand/v2"
	"os"
	"strings"
	"time"
)

type seedConfig struct {
	numServers             int
	numUsers               int
	minCategoriesPerServer int
	maxCategoriesPerServer int
	minChannelsPerCategory int
	maxChannelsPerCategory int
	heavyChannelMessages   int
	lightChannelMessages   int
	seedPrefix             string
}

func parseFlags() seedConfig {
	cfg := seedConfig{}
	flag.IntVar(&cfg.numServers, "servers", 15, "Number of servers to create")
	flag.IntVar(&cfg.numUsers, "users", 300, "Number of users to create")
	flag.IntVar(&cfg.minCategoriesPerServer, "minCategories", 2, "Minimum categories per server")
	flag.IntVar(&cfg.maxCategoriesPerServer, "maxCategories", 5, "Maximum categories per server")
	flag.IntVar(&cfg.minChannelsPerCategory, "minChannels", 2, "Minimum channels per category")
	flag.IntVar(&cfg.maxChannelsPerCategory, "maxChannels", 5, "Maximum channels per category")
	flag.IntVar(&cfg.heavyChannelMessages, "messages", 5000, "Number of messages to put in one heavy channel")
	flag.IntVar(&cfg.lightChannelMessages, "lightMessages", 500, "Number of messages to put in other channels (0 to skip)")
	flag.StringVar(&cfg.seedPrefix, "prefix", "seed", "Unique prefix to avoid unique constraint collisions")
	flag.Parse()

	if cfg.maxCategoriesPerServer < cfg.minCategoriesPerServer {
		cfg.maxCategoriesPerServer = cfg.minCategoriesPerServer
	}
	if cfg.maxChannelsPerCategory < cfg.minChannelsPerCategory {
		cfg.maxChannelsPerCategory = cfg.minChannelsPerCategory
	}
	return cfg
}

func randomSentence() string {
	words := []string{
		"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit",
		"sed", "do", "eiusmod", "tempor", "incididunt", "ut", "labore", "et", "dolore",
		"magna", "aliqua", "ut", "enim", "ad", "minim", "veniam", "quis", "nostrud",
		"exercitation", "ullamco", "laboris", "nisi", "ut", "aliquip", "ex", "ea",
		"commodo", "consequat", "duis", "aute", "irure", "dolor", "in", "reprehenderit",
		"in", "voluptate", "velit", "esse", "cillum", "dolore", "eu", "fugiat",
	}
	n := 6 + random.IntN(14)
	var b strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteString(" ")
		}
		w := words[random.IntN(len(words))]
		if i == 0 {
			if len(w) > 0 {
				w = strings.ToUpper(w[:1]) + w[1:]
			}
		}
		b.WriteString(w)
	}
	b.WriteString(".")
	return b.String()
}

func makeTipTapDoc(text string) json.RawMessage {
	content := map[string]any{
		"type": "doc",
		"content": []any{
			map[string]any{
				"type": "paragraph",
				"content": []any{
					map[string]any{
						"type": "text",
						"text": text,
					},
				},
			},
		},
	}
	b, _ := json.Marshal(content)
	return b
}

func main() {
	cfg := parseFlags()

	ctx := context.Background()
	db := database.New()
	defer db.Close()

	fmt.Println("Seeding database...")

	// 1) Users
	fmt.Println("Creating users...")
	userIDs := make([]string, 0, cfg.numUsers)
	for i := 0; i < cfg.numUsers; i++ {
		email := fmt.Sprintf("%s-user-%d-%d@example.com", cfg.seedPrefix, time.Now().Unix(), i)
		username := fmt.Sprintf("%s_user_%d_%d", cfg.seedPrefix, time.Now().Unix(), i)
		displayName := fmt.Sprintf("User %d", i+1)

		hashed, err := crypto.HashPassword("password1234")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to hash password: %v\n", err)
			os.Exit(1)
		}

		u, err := db.CreateUser(ctx, &types.SignUpParams{
			Email:       email,
			Username:    username,
			DisplayName: displayName,
			Password:    hashed,
		})
		if err != nil {
			// Skip duplicates gracefully
			fmt.Fprintf(os.Stderr, "CreateUser error: %v\n", err)
			continue
		}
		userIDs = append(userIDs, u.ID)
	}
	if len(userIDs) == 0 {
		fmt.Fprintln(os.Stderr, "no users created, aborting")
		os.Exit(1)
	}

	// 2) Servers, categories, channels
	fmt.Println("Creating servers, categories, and channels...")
	type channelRef struct {
		serverID  string
		channelID string
		memberIDs []string
	}
	var channels []channelRef

	for s := 0; s < cfg.numServers; s++ {
		owner := userIDs[random.IntN(len(userIDs))]
		name := fmt.Sprintf("%s Server %d", cfg.seedPrefix, s+1)

		desc, _ := json.Marshal(map[string]any{"ops": []any{}})
		avatars := []string{
			"https://d2vxg81co3irvv.cloudfront.net/h5l2kgkgl5388cn02q4qfxsp.webp",
			"https://d2vxg81co3irvv.cloudfront.net/avlu2a8o91gpvalk6gmgjx5r-avatar-dmplh84i8bft0ww7g60qrt15.webp",
			"https://d2vxg81co3irvv.cloudfront.net/wn3yuwwghg7efv5an78bojyc.webp",
			"https://d2vxg81co3irvv.cloudfront.net/ru3piutctxqesrv1xom8d8nn.webp",
			"https://d2vxg81co3irvv.cloudfront.net/rx38ak222ydoxb5kevww01ce-avatar-nf3531et6izcz2k6uyk5nm84.webp",
			"https://d2vxg81co3irvv.cloudfront.net/qilvsvq39rbmnybr4s24e86d-avatar-xtzdeamrkmtmh11odhtb9p2y.webp",
			"https://d2vxg81co3irvv.cloudfront.net/vlxth3gbqa3bmcsysy581pzy-avatar-frh262c7w15ntm8jia4jjrud.webp",
			"https://d2vxg81co3irvv.cloudfront.net/jnclk7c4el73kaqc726qlw7c-avatar-u0090mlv8r3zu5d850unedvu.webp",
			"https://d2vxg81co3irvv.cloudfront.net/u8b6rr2uz0xn3ubx41o6o6wx.webp",
			"https://d2vxg81co3irvv.cloudfront.net/f6983q524yum1ntx468j5gzu-avatar-mpkv59uvddg7jjf26gbigoas.webp",
			"https://d2vxg81co3irvv.cloudfront.net/qvn4bem4ms537cnn3kjyyv0l-avatar-rmk4i0evo7ls9v75jllzeuw0.webp",
		}
		avatar := avatars[random.IntN(len(avatars))]
		crop := types.Crop{X: 0, Y: 0, Width: 256, Height: 256}
		srv, err := db.CreateServer(ctx, owner, &types.CreateServerParams{
			Name:        name,
			Description: desc,
			Public:      true,
			Crop:        crop,
			Position:    s,
		}, &avatar)
		if err != nil {
			fmt.Fprintf(os.Stderr, "CreateServer error: %v\n", err)
			continue
		}

		// Join a random subset of users to this server
		numMembers := 30 + random.IntN(max(1, min(cfg.numUsers/2, 120)))
		members := make([]string, 0, numMembers)
		used := map[string]struct{}{}
		for i := 0; i < numMembers; i++ {
			uid := userIDs[random.IntN(len(userIDs))]
			if _, ok := used[uid]; ok {
				continue
			}
			if _, err := db.JoinServer(ctx, srv.ID, uid, i); err == nil {
				members = append(members, uid)
				used[uid] = struct{}{}
			}
		}
		if len(members) == 0 {
			members = []string{owner}
		}

		// Categories
		cats := cfg.minCategoriesPerServer + random.IntN(cfg.maxCategoriesPerServer-cfg.minCategoriesPerServer+1)
		categoryIDs := make([]string, 0, cats)
		for c := 0; c < cats; c++ {
			cat, err := db.CreateCategory(ctx, &types.CreateCategoryParams{
				ServerID: srv.ID,
				Name:     fmt.Sprintf("category-%d", c+1),
				Position: c,
				Users:    nil,
				Roles:    nil,
				E2EE:     false,
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "CreateCategory error: %v\n", err)
				continue
			}
			categoryIDs = append(categoryIDs, cat.ID)
		}
		if len(categoryIDs) == 0 {
			// Ensure at least one category
			cat, err := db.CreateCategory(ctx, &types.CreateCategoryParams{
				ServerID: srv.ID,
				Name:     "general",
				Position: 0,
			})
			if err == nil {
				categoryIDs = append(categoryIDs, cat.ID)
			}
		}

		// Channels
		for _, catID := range categoryIDs {
			chs := cfg.minChannelsPerCategory + random.IntN(cfg.maxChannelsPerCategory-cfg.minChannelsPerCategory+1)
			for cc := 0; cc < chs; cc++ {
				ch, err := db.CreateChannel(ctx, &types.CreateChannelParams{
					Position:    cc,
					CategoryID:  catID,
					ServerID:    srv.ID,
					Name:        fmt.Sprintf("chat-%d", cc+1),
					Description: "",
					Users:       nil,
					Roles:       nil,
					E2EE:        false,
					Type:        "textual",
				})
				if err != nil {
					fmt.Fprintf(os.Stderr, "CreateChannel error: %v\n", err)
					continue
				}
				channels = append(channels, channelRef{serverID: srv.ID, channelID: ch.ID, memberIDs: members})
			}
		}
	}

	if len(channels) == 0 {
		fmt.Fprintln(os.Stderr, "no channels created, aborting")
		os.Exit(1)
	}

	// 3) Messages
	fmt.Println("Creating messages...")
	// Pick one heavy channel
	heavy := channels[random.IntN(len(channels))]
	for i := 0; i < cfg.heavyChannelMessages; i++ {
		author := heavy.memberIDs[random.IntN(len(heavy.memberIDs))]
		body := &types.CreateMessageParams{
			ServerID:         heavy.serverID,
			ChannelID:        heavy.channelID,
			Content:          makeTipTapDoc(randomSentence()),
			Everyone:         false,
			MentionsUsers:    nil,
			MentionsRoles:    nil,
			MentionsChannels: nil,
			Attachments:      json.RawMessage("[]"),
		}
		if _, err := db.CreateMessage(ctx, author, body); err != nil {
			// keep going for throughput
			if i%100 == 0 {
				fmt.Fprintf(os.Stderr, "CreateMessage (heavy) error at %d: %v\n", i, err)
			}
		}
		if (i+1)%1000 == 0 {
			fmt.Printf("  heavy channel: %d/%d messages inserted\n", i+1, cfg.heavyChannelMessages)
		}
	}

	if cfg.lightChannelMessages > 0 {
		for _, ch := range channels {
			if ch.channelID == heavy.channelID {
				continue
			}
			n := cfg.lightChannelMessages / 2
			if n <= 0 {
				n = 1
			}
			n += random.IntN(max(1, cfg.lightChannelMessages-n+1))
			for i := 0; i < n; i++ {
				author := ch.memberIDs[random.IntN(len(ch.memberIDs))]
				body := &types.CreateMessageParams{
					ServerID:         ch.serverID,
					ChannelID:        ch.channelID,
					Content:          makeTipTapDoc(randomSentence()),
					Everyone:         false,
					MentionsUsers:    nil,
					MentionsRoles:    nil,
					MentionsChannels: nil,
					Attachments:      json.RawMessage("[]"),
				}
				_, _ = db.CreateMessage(ctx, author, body)
			}
		}
	}

	fmt.Println("Seeding complete.")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
