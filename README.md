Hosting: Google cloud run

Backend: Go (https://freshman.tech/web-development-with-go/) (https://www.youtube.com/watch?v=W5b64DXeP0o&list=PLzUGFf4GhXBL4GHXVcMMvzgtO8-WEJIoY)

Frontend: HTMX (https://htmx.org/examples)

Deploying:
  run `gcloud run deploy --source .` in source code directory
  choose `us-central1` for the region

TODO:
- restructure go packages in some intentional way
  - https://medium.com/sellerapp/golang-project-structuring-ben-johnson-way-2a11035f94bc


Docker:
  app:
    `docker run -dp 127.0.0.1:3000:3000 --mount type=volume,src=todo-db,target=/etc/todos getting-started`
    build: `docker build --tag htmx-demo .` (https://docs.docker.com/language/golang/build-images/)
    run: `docker run --publish 6969:6969 htmx-demo` (https://docs.docker.com/language/golang/run-containers/)
  sql: (https://docs.docker.com/get-started/07_multi_container/)
    create network: `docker network create htmx-demo-db`
    start & attach sql: `docker run -d --network htmx-demo-db --network-alias mysql -v htmx-demo-db:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=secret -e MYSQL_DATABASE=demo-db mysql:8.0`

  usage:
    Don't do anything `docker run` in the terminal, this is handled by docker compose


SQL db:
  Instance ID: 'htmx-demo-db'
  password: password