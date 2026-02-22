# R2 Manager 機能ギャップ分析

他のS3互換ストレージ管理ツール（AWS S3 Console、S3 Browser、Cyberduck、MinIO Console、CloudBerry Explorer、Filestash、TagSpaces、R2 Client等）と比較した、R2 Managerの不足機能の調査結果。

## 調査日

2026-02-22

## 現在実装済みの機能

- バケット一覧表示（インメモリキャッシュ付き、60分TTL）
- オブジェクトの階層的ブラウジング（ページネーション、プレフィックス/デリミタ対応）
- ファイルアップロード（複数ファイル、ドラッグ&ドロップ、SSE進捗追跡、上書き制御）
- ディレクトリ作成（マーカーオブジェクト方式）
- ファイルダウンロード / コンテンツストリーミング
- 画像プレビュー（jpg, png, gif, webp, svg, bmp, ico）
- ファイルメタデータ表示（名前、キー、サイズ、更新日時、ETag）
- パブリックURL生成（バケットごとに静的URL設定）
- 内部URL / パブリックURLのコピー
- キャッシュ管理（インメモリ + ディスクの2層キャッシュ、ETag整合性検証付き）
- バケット設定管理（パブリックURL設定、SQLite永続化）

---

## 不足機能の分析

### 1. 基本的なファイル操作（重要度：高）

現在最も大きなギャップ。全ての競合ツールが対応している基本操作が欠けている。

| 不足機能 | 対応ツール例 | 実装難度 | 備考 |
|---------|------------|---------|------|
| **ファイル削除** | 全ツール | 低 | S3 DeleteObject API。最も基本的な操作 |
| **ファイルリネーム** | S3 Browser, Cyberduck, CloudBerry, Filestash | 中 | S3にはネイティブのrenameがないため、CopyObject + DeleteObject で実現 |
| **ファイル移動** | AWS Console, Cyberduck, Filestash, MinIO | 中 | リネームと同様、copy + delete パターン |
| **ファイルコピー** | AWS Console, S3 Browser, CloudBerry | 低〜中 | CopyObject API。同一バケット内/バケット間の両方 |
| **一括選択・一括操作** | AWS Console, MinIO, S3 Browser, Cyberduck | 中 | 複数ファイル選択 → 一括削除・移動・コピー。UIとAPIの両面で対応必要 |

### 2. ファイル共有（重要度：高）

| 不足機能 | 対応ツール例 | 実装難度 | 備考 |
|---------|------------|---------|------|
| **署名付きURL（Presigned URL）** | AWS Console (12h), Cyberduck, S3 Browser, CloudBerry, s3cmd | 低〜中 | R2のS3互換APIでPresigned URLをサポート。期限付き共有リンクの生成。個人用途でも外部共有に便利 |

### 3. 検索・フィルタリング・ソート（重要度：高）

ファイル数が増えると閲覧・管理が困難になるため重要。

| 不足機能 | 対応ツール例 | 実装難度 | 備考 |
|---------|------------|---------|------|
| **オブジェクト名検索** | MinIO Console, S3 Browser, TagSpaces, R2 Uploader | 低〜中 | プレフィックスベースのフィルタリング、部分一致検索 |
| **ファイルタイプフィルタ** | R2 Client, TagSpaces, Cyberduck | 低 | 拡張子・MIMEタイプでの絞り込み |
| **カラムソート** | ほぼ全ツール | 低 | 名前・サイズ・更新日時等でのソート切り替え |

### 4. プレビュー機能の拡充（重要度：中）

現在は画像のみ。他ツールは多様なファイルタイプのプレビューを提供。

| 不足機能 | 対応ツール例 | 実装難度 | 備考 |
|---------|------------|---------|------|
| **テキストファイルプレビュー** | Filestash, TagSpaces, CloudBerry | 低 | .txt, .json, .md, .csv, .yaml等。ブラウザ内で表示可能 |
| **PDFプレビュー** | TagSpaces, Filestash | 中 | ブラウザ内蔵のPDFビューアまたはpdf.jsで実現可能 |
| **動画/音声プレビュー** | TagSpaces, Filestash | 低〜中 | HTML5 video/audioタグでブラウザネイティブ対応 |
| **コードハイライト** | Filestash, TagSpaces | 低 | シンタックスハイライトライブラリ（Prism.js等）利用 |
| **Markdownレンダリング** | Filestash, TagSpaces | 低 | marked/react-markdown等で実装可能 |

### 5. ストレージ分析・統計（重要度：中）

| 不足機能 | 対応ツール例 | 実装難度 | 備考 |
|---------|------------|---------|------|
| **バケット使用量表示** | AWS Console (Storage Lens), MinIO Console, R2 Client, R2 Dashboard | 中 | バケットごとの総容量・オブジェクト数 |
| **ファイルタイプ分布** | R2 Client | 中〜高 | ファイル種別の円グラフ等。全オブジェクトの走査が必要 |
| **アップロード推移** | R2 Client | 中 | 日別のアップロードトレンド。ローカルDB記録で実現可能 |

### 6. バケット管理の拡充（重要度：中）

| 不足機能 | 対応ツール例 | 実装難度 | 備考 |
|---------|------------|---------|------|
| **バケット作成** | AWS Console, MinIO Console, S3 Browser, s3cmd | 低 | CreateBucket API |
| **バケット削除** | AWS Console, MinIO Console, S3 Browser | 低〜中 | DeleteBucket API。空でない場合の処理が必要 |
| **CORS設定** | AWS Console, MinIO Console, S3 Browser | 中 | PutBucketCors API。R2対応状況の確認が必要 |

### 7. メタデータ管理（重要度：中〜低）

| 不足機能 | 対応ツール例 | 実装難度 | 備考 |
|---------|------------|---------|------|
| **カスタムメタデータ表示** | Cyberduck, S3 Browser (Pro), CloudBerry | 低 | HeadObject APIでContent-Type, Cache-Control等を取得 |
| **メタデータ編集** | Cyberduck, S3 Browser (Pro) | 中 | CopyObject（自分自身へ）でメタデータ上書き |
| **オブジェクトタグ** | AWS Console, MinIO Console | 中 | GetObjectTagging / PutObjectTagging API |

### 8. その他の機能（重要度：低〜中）

| 不足機能 | 対応ツール例 | 実装難度 | 備考 |
|---------|------------|---------|------|
| **フォルダごとアップロード** | AWS Console, Cyberduck, CloudBerry | 中 | webkitdirectoryを使ったディレクトリ一括アップロード |
| **ドラッグ&ドロップ移動（フォルダ間）** | ForkLift, Cyberduck | 高 | UI上でのD&D移動はcopy+delete操作に変換 |
| **テキストファイル直接編集** | Filestash, CloudBerry, TagSpaces | 中 | インプレースエディタ → PutObject で保存 |
| **重複ファイル検出** | R2 Client | 高 | ETagやハッシュベースの検出。全オブジェクト走査必要 |
| **キーボードショートカット** | S3 Browser, Cyberduck | 低〜中 | Delete, Ctrl+C/V等の基本操作 |

---

## R2固有の制約

以下の機能はCloudflare R2自体が未サポートまたは制限があるため、実装対象外とする。

- **バージョニング** — R2は未サポート
- **ライフサイクルポリシー** — R2は未サポート（オブジェクト有効期限のみ対応）
- **暗号化設定** — R2はデフォルトで暗号化。ユーザー管理キーは未サポート
- **ストレージクラス管理** — R2は単一ストレージクラス（Infrequent Access追加予定）
- **レプリケーション** — R2は未サポート
- **Object Lock / WORM** — R2は未サポート

---

## 推奨実装優先順位

個人利用のセルフホスト用途という本プロジェクトの性質を考慮した優先順位。

### Phase 1（基本操作の完成）
1. **ファイル削除**（単一 + 一括）
2. **カラムソート**（名前・サイズ・更新日時）
3. **オブジェクト名検索/フィルタ**

### Phase 2（ファイル管理の充実）
4. **ファイルリネーム**（copy + delete）
5. **ファイル移動**（copy + delete）
6. **ファイルコピー**
7. **署名付きURL（Presigned URL）生成**

### Phase 3（プレビュー拡充）
8. **テキスト/コード/Markdownプレビュー**
9. **PDF プレビュー**
10. **動画/音声プレビュー**

### Phase 4（管理機能の強化）
11. **バケット使用量統計**
12. **バケット作成/削除**
13. **カスタムメタデータ表示**
14. **フォルダごとアップロード**

---

## 調査対象ツール

| ツール名 | 種類 | プラットフォーム | ライセンス |
|---------|------|--------------|----------|
| AWS S3 Console | Webコンソール | Web | 無料（AWSアカウント必要） |
| S3 Browser | デスクトップGUI | Windows | 無料（個人）/ $29.95（Pro） |
| Cyberduck | デスクトップGUI | Mac/Windows | 無料（OSS） |
| MinIO Console | Webコンソール | Web | 無料（コミュニティ版、管理機能制限あり） |
| CloudBerry Explorer (MSP360) | デスクトップGUI | Windows/Mac | 無料 / $39.99（Pro） |
| s3cmd | CLI | 全OS | 無料（OSS） |
| Filestash | セルフホストWeb | Web | 無料（OSS） |
| TagSpaces | デスクトップ + PWA | 全OS | 無料（Lite）/ 有料（Pro） |
| R2 Client | デスクトップGUI | 全OS | 有料 |
| R2 Uploader | デスクトップGUI | 全OS | 無料（OSS） |
| Cloudflare R2 Dashboard | Webコンソール | Web | 無料 |
