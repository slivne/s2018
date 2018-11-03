package main

import (
	"fmt"
	"log"
	"time"
	"math/rand"

	"github.com/gocql/gocql"
)

var users = make([]string, 100)

func generate_users(){
	for i := 0; i < len(users); i++ {
		users[i] = fmt.Sprintf("user%d",i)
	}
}

func get_followers(user string) []string {
	var id int
	fmt.Sscanf(user,"user%d",&id)
	followers := make([]string, 5)
	for i := 0; i < 5; i++ {
		followers[i] = users[(id+i+1)%len(users)]
	}
	return  followers
}


func insert_tweet(session *gocql.Session, user string, tweet_id gocql.UUID, tweet_time gocql.UUID, tweet_txt string){
	if err := session.Query("INSERT INTO tweets (user, tweet_id, time, text) VALUES ( ?, ?, ?, ?)",
		user, tweet_id, tweet_time, tweet_txt).Exec(); err != nil {
			log.Fatal(err)
	}


	for _, follower := range get_followers(user) {
		liked := false
		if (rand.Intn(1000) <= 1) {
			liked = true
		}
		if err := session.Query("INSERT INTO timeline (user, tweet_id, time, author, text, liked) VALUES ( ?, ?, ?, ?, ?, ?)",
			follower, tweet_id, tweet_time, user, tweet_txt, liked).Exec(); err != nil {
				log.Fatal(err)
		}
	}
}

func get_timeline(session *gocql.Session, user string, filter_liked bool) {
	var tweet_id gocql.UUID
	var author string
	var text string
	var liked bool

	if !filter_liked {
		query := "SELECT tweet_id, author, text, liked FROM timeline WHERE user = ? limit 50"
		iter := session.Query(query, user).Iter()
		for iter.Scan(&tweet_id, &author, &text, &liked) {}
		if err := iter.Close(); err != nil {
			log.Fatal(err)
		}
	} else {
		query := "SELECT tweet_id, author, text, liked FROM timeline_liked WHERE user = ? and liked = ? limit 50"
		iter := session.Query(query, user, true).Iter()
		for iter.Scan(&tweet_id, &author, &text, &liked) {}
		if err := iter.Close(); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	// connect to the cluster
        cluster := gocql.NewCluster("172.17.0.2", "172.17.0.3", "172.17.0.4")

	cluster.Keyspace = "scylla_demo"
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 5000000000
        cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy());
	session, _ := cluster.CreateSession()
	defer session.Close()

	generate_users()

	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	rate := time.Second / 20
	throttle := time.Tick(rate)
	for (true) {
		<-throttle
		for count := 0; count < 10; count++ {
			user := users[random.Intn(len(users))]
			for msg := 0 ; msg < 1; msg++ {
				insert_tweet(session, user, gocql.TimeUUID(), gocql.TimeUUID(), fmt.Sprintf("msg_%s_%d",user,msg))
			}
			get_timeline(session, user, random.Intn(10) < 5)
		}
	}
}
