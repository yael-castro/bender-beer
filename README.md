# bender-beer
Test MS Beer

###### How to run
```
sudo docker rm -f bender-beer-container
sudo docker-compose build
sudo docker run -p 80:8080 --name bender-beer-container --env-file env.env bender-beer-image
```
