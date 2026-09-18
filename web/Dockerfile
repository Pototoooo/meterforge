FROM node:24-alpine AS build
WORKDIR /workspace
ENV CI=true

COPY web/package.json web/pnpm-lock.yaml web/pnpm-workspace.yaml ./web/
RUN corepack enable && cd web && pnpm install --frozen-lockfile

COPY api/client/javascript ./api/client/javascript
COPY web ./web
RUN cd web && pnpm build

FROM nginx:1.27-alpine
COPY web/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /workspace/web/dist /usr/share/nginx/html
EXPOSE 8080
HEALTHCHECK --interval=5s --timeout=3s --retries=20 CMD wget --spider http://127.0.0.1:8080/healthz || exit 1
