sudo docker run --rm --name myubuntu --network host ubuntu bash -c "apt-get update && apt-get install -y curl && curl http://localhost:8080"

docker run -d --name mynginx -p 8080:80 nginx

cat ebpf_exporter.log |grep CNetwork_event  | wc -l