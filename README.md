Ticket Expert
Step to run the project:
- Add environment variables in .env
- install dependencies in go.mod
- Run docker-compose.yaml, to start Redis, Postgres, and Simple upload server. 
- Run main.go


export $(cat .env | xargs) 