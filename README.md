Hosting: Google cloud run

Backend: Go (https://freshman.tech/web-development-with-go/) (https://www.youtube.com/watch?v=W5b64DXeP0o&list=PLzUGFf4GhXBL4GHXVcMMvzgtO8-WEJIoY)

Frontend: HTMX (https://htmx.org/examples)

Deploying:
  run `gcloud run deploy --source .` in source code directory
  choose `us-central1` for the region

TODO:
- restructure go packages in some intentional way
  - https://medium.com/sellerapp/golang-project-structuring-ben-johnson-way-2a11035f94bc


usage:
  Don't do anything like `docker run` in the terminal, this is handled by docker compose
  Do: `docker compose down`
      `docker compose up`

How to launch this baddie:
  `docker build --tag htmx-demo .`
  `docker compose up`
Alternatively, just run `air` (lazy)

Local instance at `localhost:6969`

SQL db:
  Instance ID: 'htmx-demo-db'
  password: password


templ: Templating
`~/go/bin/templ generate`
