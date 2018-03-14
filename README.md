# s18-fish-findr

By: Jason Alexander, Mitchell Mckenzie, Zachary Whitworth, Carson Sallis

This is for a 4910 project. The purpose of the mobile app is to help people who are fishing alone or want to fish with new people, find new buddies to fish with. It will accept personal information such as likes, dislikes, type of fishing you enjoy, etc. Once this info is compiled into your profile then it will match you with similar fisherman on the app to get in contact with. You can ever search for fisherman with interest that you specify.

# Getting Started For Development

## Windows
Download the [docker container engine](https://store.docker.com/editions/community/docker-ce-desktop-windows).

If you would like to ensure that your docker engine is working well enough to run the DB, visit this [site](https://docs.docker.com/docker-for-windows/) and go up to step 3. 

## Mac
Download the [docker container engine](https://store.docker.com/editions/community/docker-ce-desktop-mac).

If you would like to ensure that your docker engine is working well enough to run the DB, visit this [site](https://docs.docker.com/docker-for-mac/) and go up to step 3. 

## Linux
Download the [docker container engine](https://store.docker.com/editions/community/docker-ce-server-ubuntu).

If you would like to ensure that your docker engine is working well enough to run the DB, visit this [site](https://docs.docker.com/get-started/part2/) and go until you are comfortable. 

## Everyone
Download [golang](https://golang.org/dl/). Follow the install instructions [here](https://golang.org/doc/install).

Clone this repo into your gopath at `$GOPATH/src/github.com/Clemson-CPSC-4910`

Once inside the repo, follow these commands to start and use the app.

    docker-compose build
    docker-compose up

You can test that the api has started by visiting `http://localhost:8002/` to see the webapp.

If you just want to start the database then you can run `docker-compose up postgres` after you have built.

## More information on golang and development in the README inside /go

# WebApp
Developnent for the webapp will take place in the /go/webapp directory. Place everything to do with the web app including all html/js/css in there.

# Testing

## Go/Postgres Test
These test are tagged with the tag `db`. In order to run these test you should start the postgres database. 

    docker-compose build
    docker-compose up postgres

Then you can run the test like so: 

    cd fisherman-finder/go/postgres
    go test -tags=db