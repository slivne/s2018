package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gocql/gocql"
)

var (
	hosts       = flag.String("hosts", "", "comma separated list of hosts to connect to")
	users       = make([]string, 1000)
)

func main() {

	flag.Parse()

	if *hosts == "" {
		flag.Usage()
		return
	}

	// connect to the cluster
	cluster := gocql.NewCluster(strings.Split(*hosts, ",")...)
	cluster.Keyspace = "scylla_demo"
        cluster.Consistency = gocql.Quorum
        cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy());
        cluster.Timeout = 5000000000
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	ctx, cancel := context.WithCancel(context.Background())
	quitSignal(cancel)
	generateUsers()

	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	rate := time.Second / 20
	throttle := time.NewTicker(rate)
	for {
		select {
		case <-throttle.C:
			for count := 0; count < 10; count++ {
				user := users[random.Intn(len(users))]
                                insertTweet(session, user, gocql.TimeUUID(), gocql.TimeUUID(), fmt.Sprintf("msg_%s_%s",user,time.Now().String()))
				getTimeline(session, user, random.Intn(10) < 1)
			}
		case <-ctx.Done():
			throttle.Stop()
			return
		}
	}
}

func generateUsers() {
	for i := 0; i < len(users); i++ {
		users[i] = fmt.Sprintf("user%d", i)
	}
}

func getFollowers(user string) []string {
	var id int
	fmt.Sscanf(user, "user%d", &id)
	followers := make([]string, 5)
	for i := 0; i < 5; i++ {
		followers[i] = users[(id+i+1)%len(users)]
	}
	return followers
}

func insertTweet(session *gocql.Session, user string, tweetID gocql.UUID, tweetTime gocql.UUID, tweetTxt string) {
	if err := session.Query("INSERT INTO tweets (user, tweet_id, time, text) VALUES (?, ?, ?, ?)",
		user, tweetID, tweetTime, tweetTxt).Exec(); err != nil {
		log.Fatal(err)
	}

	for _, follower := range getFollowers(user) {
		liked := false
		if rand.Intn(1000) < 1 {
			liked = true
		}
		if err := session.Query("INSERT INTO timeline (user, tweet_id, time, author, text, liked) VALUES (?, ?, ?, ?, ?, ?)",
			follower, tweetID, tweetTime, user, tweetTxt, liked).Exec(); err != nil {
			log.Fatal(err)
		}
	}
}

func getTimeline(session *gocql.Session, user string, filter_liked bool) {
	var tweetID gocql.UUID
	var author string
	var text string

	if (!filter_liked) {
		iter := session.Query("SELECT tweet_id, author, text FROM timeline WHERE user = ? limit 50", user).Iter()
		for iter.Scan(&tweetID, &author, &text) {
		}
		if err := iter.Close(); err != nil {
			log.Fatal(err)
		}
	} else {
		iter := session.Query("SELECT tweet_id, author, text FROM timeline_liked WHERE user = ? and liked = ? limit 50", user, true).Iter()
		for iter.Scan(&tweetID, &author, &text) {
		}
		if err := iter.Close(); err != nil {
			log.Fatal(err)
		}
	}
}

func quitSignal(cancel context.CancelFunc) chan bool {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-sigs:
			cancel()
		}
	}()

	return done
}
