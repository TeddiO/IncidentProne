# IncidentProne
A small ISMS / ITIL style system designed for explicitly tracking and handling incidents (PoC for a Hackathon).

While we have Jira, it's very easy for incidents to get lost in there alongside that not everyone has a Jira account! IncidentProne is a proof of concept system that would allow anyone with the appropriate access (which is currently missing 😂) to log in and report / view / resolve incidents. Additionally if we did want this data replicated to a system like Jira, this is easily implemented and tickets could be raised based on priority. 

Also missing (due to time constraints) would be straightforward reporting of how many incidents we've had, how many are still open and so on. 

## Tech used
Application is powered by Go which is backed by a Postgres database. We use docker to simplify setup and running the example for presentation purposes.
Any and all assets are actually baked in to the executable (totalling 15MB in total) and could be ran on any system without any requirements of any other dependencies. 

### Setup
You'll want to fire up the initial compose file to fetch postgres

```bash 
docker-compose -f ./docker-compose.yml up -d
```

Then execute the database-setup.sh file if you haven't got any data in the database
```bash 
./database-setup.sh
```

And then finally.
```bash 
docker-compose -f ./docker-compose.yml up -d --force-recreate
```

The application should be available on `localhost:8080` after a few moments!
