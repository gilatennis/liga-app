# Season Setup

1. Assign teams to groups on [gsheet](https://docs.google.com/spreadsheets/d/1OT-r1U4_xeRmg53Wj62ujbfto0aaURXynyu80ka2vOQ/edit?gid=0#gid=0)
2. generate pool matches
  ```bash
  $ go run cmd/seed/main.go generate-pool-matches --league-uuid <league-uuid> --home-and-away <true/false>
  ```
