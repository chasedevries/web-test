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
  build: `docker build --tag htmx-demo .` (https://docs.docker.com/language/golang/build-images/)
  run: `docker run --publish 6969:6969 htmx-demo` (https://docs.docker.com/language/golang/run-containers/)

SQL db:
  Instance ID: 'htmx-demo-db'
  password: password