## 1.0.0 (2025-12-21)


### 🚀 Features

* add Air configuration for hot-reload ([0d1b57f](https://github.com/daniel-caso-github/realtime-alert-system/commit/0d1b57f945fec07f86f20421a29beb3bdb30a0e1))
* add alert and auth handlers with validation helpers ([41a4122](https://github.com/daniel-caso-github/realtime-alert-system/commit/41a4122471cbd4078a7edc8b2363d46dfa256041))
* add alert event producer for async event publishing ([39b081b](https://github.com/daniel-caso-github/realtime-alert-system/commit/39b081b591c964d3613056da57f18184b363e69d))
* add AlertManager integration with webhook handler ([aa4bc2a](https://github.com/daniel-caso-github/realtime-alert-system/commit/aa4bc2aa0f410ec8d0cef4caffef683a49f367ee))
* add application layer DTOs and services (auth, alert) ([399015d](https://github.com/daniel-caso-github/realtime-alert-system/commit/399015dc41c42aa1325a8530206e6e9d9a924c98))
* add CI/CD pipeline with GitHub Actions ([082a716](https://github.com/daniel-caso-github/realtime-alert-system/commit/082a716b68c2e0deade228919f21b0cc9cf56043))
* add circuit breaker for external services ([1001e9b](https://github.com/daniel-caso-github/realtime-alert-system/commit/1001e9bf1febe99e772a2549dcdd1c519c3239e1))
* add configuration management with Viper ([7ee0cdb](https://github.com/daniel-caso-github/realtime-alert-system/commit/7ee0cdbed8c49ef36dde95396232ab5156709c10))
* add database initialization scripts ([b6edfe1](https://github.com/daniel-caso-github/realtime-alert-system/commit/b6edfe14f47dd5f7c2a28e5d6b61b866b4b3bc58))
* add database migrations with golang-migrate ([b49ab27](https://github.com/daniel-caso-github/realtime-alert-system/commit/b49ab273832a63302f3bcb6f0e9ae07f3d17092d))
* add dead letter queue processor with admin endpoints ([41a0ad5](https://github.com/daniel-caso-github/realtime-alert-system/commit/41a0ad53e4fb076b4741b65d24754773b0eab870))
* add development utility scripts ([92e88f7](https://github.com/daniel-caso-github/realtime-alert-system/commit/92e88f792a454151ccad99cb4167eb39b402b478))
* add event consumers with consumer groups ([4d05dda](https://github.com/daniel-caso-github/realtime-alert-system/commit/4d05ddaef8c6be1d1c115c6be2ee9224c250ec37))
* add GitHub Actions CI workflow ([3bd3804](https://github.com/daniel-caso-github/realtime-alert-system/commit/3bd38045ccb150bc65e17d9eb50e81ab7de97e43))
* add Grafana with provisioned dashboards ([53ad711](https://github.com/daniel-caso-github/realtime-alert-system/commit/53ad711dc203724748d818c95901dfccea6e5dc6))
* add health check endpoints and HTTP router ([d66e0a6](https://github.com/daniel-caso-github/realtime-alert-system/commit/d66e0a6ef5720cec834c5ff9c465e8e2079b2744))
* add JWT authentication and role middleware ([85ca46c](https://github.com/daniel-caso-github/realtime-alert-system/commit/85ca46c17d27c616d80731d38a1c440e206e4e26))
* add Prometheus metrics instrumentation ([23b291d](https://github.com/daniel-caso-github/realtime-alert-system/commit/23b291d80a9aa51e3aa74a2cd2e62a92f4c35596))
* add Prometheus server with exporters for Redis and PostgreSQL ([789c7c8](https://github.com/daniel-caso-github/realtime-alert-system/commit/789c7c8eb77365ab25c9803be9c5d488e6a61588))
* add rate limiting middleware with Redis backend ([b3e4e85](https://github.com/daniel-caso-github/realtime-alert-system/commit/b3e4e85c94565579eaaac9f5625ec9fa0a0075f1))
* add Redis cache implementation and cached user repository ([7776ac6](https://github.com/daniel-caso-github/realtime-alert-system/commit/7776ac696f25b26c89b51db02322f6ce1b5beb8c))
* add Redis Streams event bus infrastructure ([6480fbc](https://github.com/daniel-caso-github/realtime-alert-system/commit/6480fbc616461791382003966f9636ac67f04a2a))
* add repository interfaces (ports) for all entities ([ca22e4b](https://github.com/daniel-caso-github/realtime-alert-system/commit/ca22e4bdac7f28ad5dd29477f13394e84032f525))
* add retry logic with exponential backoff ([1dd8c1b](https://github.com/daniel-caso-github/realtime-alert-system/commit/1dd8c1b20a9cf020b652be9b222e3dc96f0694c5))
* add Slack notification service with rate limiting ([cb83afd](https://github.com/daniel-caso-github/realtime-alert-system/commit/cb83afd2f78e8779aa8bf33cf65b660430ba6946))
* add structured logging package with context propagation ([0286f0c](https://github.com/daniel-caso-github/realtime-alert-system/commit/0286f0c5260e07d9be583d5403fcd8adf51b9e8c))
* add value objects (Email, Password, Pagination, AlertFilter) ([d2ffacc](https://github.com/daniel-caso-github/realtime-alert-system/commit/d2ffacc51a306072a65a0f1c6c5e4a04caaa72ad))
* add WebSocket server with hub for real-time broadcasting ([a04b5bb](https://github.com/daniel-caso-github/realtime-alert-system/commit/a04b5bbf9c8f6a126759873b9fd461c57fc2fc40))
* complete Phase 2 - Docker development environment ([3e04e58](https://github.com/daniel-caso-github/realtime-alert-system/commit/3e04e58e7653f29496d92080b8b2bde22376947f))
* complete WebSocket integration with real-time alert broadcasting ([05cec5d](https://github.com/daniel-caso-github/realtime-alert-system/commit/05cec5dd7ffbf133dfc9cfafa72a56da86def86a))
* integrate distributed tracing across services ([783aa58](https://github.com/daniel-caso-github/realtime-alert-system/commit/783aa58d80170612c968e59c496d9f5199d3f690))


### 🐛 Bug Fixes

* add missing event producer call in Acknowledge method ([8df3cda](https://github.com/daniel-caso-github/realtime-alert-system/commit/8df3cda90757fee6f07a75e8ebe131a48b950e79))
* add packages write permission for GHCR ([1beb421](https://github.com/daniel-caso-github/realtime-alert-system/commit/1beb421e55ca1da9e1be49274e2877b7eb27b49c))
* add Prometheus metrics to WebSocket hub ([f5badfd](https://github.com/daniel-caso-github/realtime-alert-system/commit/f5badfd34ecfa2fbb17392a9404e5a000f2aab1f))
* generate Swagger docs during Docker build ([6c7259d](https://github.com/daniel-caso-github/realtime-alert-system/commit/6c7259d1af7c2ad1bfc6c260c6bb30e0c692a7f9))
* install golangci-lint from source for Go 1.24 support ([074079d](https://github.com/daniel-caso-github/realtime-alert-system/commit/074079d5eb8e75af93a6599eafb50f71c6548a65))
* simplify security scan to avoid CodeQL permission issues ([743516c](https://github.com/daniel-caso-github/realtime-alert-system/commit/743516c8d314bf8c75357d46abb2dd4088ebe80a))
* update golangci-lint version for Go 1.24 compatibility ([9fbfafc](https://github.com/daniel-caso-github/realtime-alert-system/commit/9fbfafc062b00d7db70e57099b0888a1ac2976e1))
* use generic type in Swagger annotation for FailedEvent ([c6b4e65](https://github.com/daniel-caso-github/realtime-alert-system/commit/c6b4e659fd168403c49f1d616a21f95eb9699098))
* use GetReqHeaders instead of deprecated VisitAll ([f9d6ecb](https://github.com/daniel-caso-github/realtime-alert-system/commit/f9d6ecbcc2178d782b6e1c814acccb68a339173b))


### 📚 Documentation

* add comprehensive README documentation ([e5a7de9](https://github.com/daniel-caso-github/realtime-alert-system/commit/e5a7de9ba0322ba05e3a2f450d5db2f230c7224f))
* add Swagger/OpenAPI documentation for REST API ([e29af82](https://github.com/daniel-caso-github/realtime-alert-system/commit/e29af821b921f3a49c967eff205ce1b6aaf8d2f4))


### ✅ Tests

* add integration tests for REST API endpoints ([67c1998](https://github.com/daniel-caso-github/realtime-alert-system/commit/67c1998b42f7f51ff41a1d447c64db7a60758e50))
* add unit tests following hexagonal architecture structure ([a16187d](https://github.com/daniel-caso-github/realtime-alert-system/commit/a16187d070de91813ccbfd682ffe0c75fce9f964))


### 🔧 Maintenance

* add docker-compose configuration ([813fa21](https://github.com/daniel-caso-github/realtime-alert-system/commit/813fa21ce7f60aed0d0852f73a47f434fcd3e6fe))
* add golangci-lint configuration and editorconfig ([171ffda](https://github.com/daniel-caso-github/realtime-alert-system/commit/171ffdadde8f9d19968a7fe484831634c1cb0e51))
* add Makefile with development commands ([98d5d47](https://github.com/daniel-caso-github/realtime-alert-system/commit/98d5d4787f71f16e8f98d8c7c9a46331d4fed3a9))
* add multi-stage Dockerfile ([74c605b](https://github.com/daniel-caso-github/realtime-alert-system/commit/74c605b98cf70d4e0e2121cbeefb6cb9fb0f0bec))
* initialize Go modules with core dependencies ([3d757be](https://github.com/daniel-caso-github/realtime-alert-system/commit/3d757bef3e08f7b9cec613ef0752772bbae96f8b))
* initialize project structure with Clean Architecture ([ad41fec](https://github.com/daniel-caso-github/realtime-alert-system/commit/ad41fec00a02d6d935c30d21f67a3653fb07d4cc))
* update .gitignore ([27acc3c](https://github.com/daniel-caso-github/realtime-alert-system/commit/27acc3c905dc3430a2c6ed2f70d33de8a88f48e0))
* update version golang ([9dfb086](https://github.com/daniel-caso-github/realtime-alert-system/commit/9dfb08672df7fa344e40ded5733684f471937db3))
