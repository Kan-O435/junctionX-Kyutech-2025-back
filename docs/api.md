# バックエンド API 概要

- ベース URL: `http://localhost:{PORT}`（デフォルト `8080`）
- フレームワーク: Gin
- ルート定義: `internal/api/routes/router.go`
- 環境変数: `PORT`, `GIN_MODE`, `DATABASE_URL`, `REDIS_URL`, `JWT_SECRET`, `SPACETRACK_USERNAME`, `SPACETRACK_PASSWORD`, `PHYSICS_TICK_RATE`, `ORBIT_PROPAGATION_STEP`

---

## 健康チェック

- メソッド: GET
- パス: `/health`
- レスポンス例:

```json
{ "status": "ok", "message": "Satellite Game Backend is running!" }
```

---

## API v1: 衛星系

### 衛星の軌道情報取得

- メソッド: GET
- パス: `/api/v1/satellite/{id}/orbit`
- 取得情報:
  - `satellite_id`: 指定 ID
  - `timestamp`: 取得時刻
  - `position`: `{ x, y, z }` (km)
  - `velocity`: `{ x, y, z }` (km/s)
  - `altitude`: 高度(km)
  - `orbital_speed`: 軌道速度(km/s)

### 衛星のステータス取得

- メソッド: GET
- パス: `/api/v1/satellite/{id}/status`
- 取得情報:
  - `position`: 位置ベクトル(km)
  - `velocity`: 速度ベクトル(km/s)
  - `fuel`: 燃料残量(kg)
  - `health`: 健康状態

### 軌道制御（マニューバ）

- メソッド: POST
- パス: `/api/v1/satellite/{id}/maneuver`
- リクエスト(JSON):

```json
{
  "player_id": "string",
  "thrust_vector": { "x": 0, "y": 0, "z": 0 },
  "duration": 10
}
```

- レスポンス例: 成功、消費燃料、想定`delta_v` 等

### 利用可能な衛星一覧

- メソッド: GET
- パス: `/api/v1/satellite/available`
- 取得情報:
  - `satellites[]`: 各衛星の`id`, `name`, `type`, `resolution`, `coverage`, `status`, `capabilities[]`

### リアルタイム映像取得

- メソッド: GET
- パス: `/api/v1/satellite/video/realtime`
- クエリ:
  - `latitude` (必須, -90〜90)
  - `longitude` (必須, -180〜180)
  - `zoom` (任意, 1〜20)
  - `required_resolution` (任意)
  - `prefer_satellite` (任意)
- 取得情報: 選択された衛星、動画 URL/サムネ、品質指標、次回更新時刻など

### 映像履歴取得

- メソッド: GET
- パス: `/api/v1/satellite/video/history`
- クエリ: `latitude`, `longitude`, `hours`(整数, 既定 24)
- 取得情報: 過去記録リスト（各レコードに`timestamp`, `video_url`, `thumbnail_url`, `satellite_id`, `quality`）

### 複数衛星同時観測

- メソッド: POST
- パス: `/api/v1/satellite/video/multi-view`
- リクエスト(JSON):

```json
{
  "latitude": 0,
  "longitude": 0,
  "zoom": 10,
  "satellite_ids": ["himawari8", "terra"]
}
```

- 取得情報: 各衛星の動画 URL、解像度、更新時刻など

### ライブストリーミング開始/停止

- メソッド: POST
- パス: `/api/v1/satellite/video/stream/start`
- リクエスト(JSON): `satellite_id`, `latitude`, `longitude`, `duration_minutes`
- メソッド: POST
- パス: `/api/v1/satellite/video/stream/stop`
- リクエスト(JSON): `stream_id`

### 衛星詳細/カバレッジ

- メソッド: GET
- パス: `/api/v1/satellite/{id}/info`
- メソッド: GET
- パス: `/api/v1/satellite/{id}/coverage`
- 取得情報: 視野範囲、現在位置、次回通過など

### 地点カバレッジ

- メソッド: GET
- パス: `/api/v1/satellite/coverage/location`
- クエリ: `latitude`, `longitude`
- 取得情報: 観測可能衛星一覧、推奨衛星など

---

## API v1: 自然災害

### アクティブ災害一覧

- メソッド: GET
- パス: `/api/v1/disaster/active`
- 取得情報: `id`, `type`, `magnitude`, `location{ latitude, longitude, country }`, `severity`, `status`

### 災害詳細

- メソッド: GET
- パス: `/api/v1/disaster/{id}`
- 取得情報: 詳細情報（震源/気象など）、衛星観測サンプル

### 災害リアルタイム映像

- メソッド: GET
- パス: `/api/v1/disaster/{id}/video`
- 取得情報: `video_streams[]`（衛星 ID、URL、タイプ、更新時刻）

---

## API v1: ミッション/デブリ

### デブリ脅威一覧

- メソッド: GET
- パス: `/api/v1/mission/debris/{id}/threats`
- 取得情報: 各デブリの`id`, `name`, `distance`, `time_to_impact`, `danger_level`

---

## 静的ファイル配信

- `/api/v1/files/satellite` → `./data/satellite_images`
- `/api/v1/files/disaster` → `./data/disaster_images`
- `/api/v1/files/thumbnails` → `./data/thumbnails`

---

## WebSocket エンドポイント（プレースホルダ）

- GET `/api/v1/ws/disaster`
- GET `/api/v1/ws/satellite/stream`
  - いずれも現状は説明 JSON を返すプレースホルダ
