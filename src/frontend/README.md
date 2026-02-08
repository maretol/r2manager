# R2 Manager - Frontend

Cloudflare R2のコンテンツを管理するためのWebアプリケーションのフロントエンドです。

## 技術スタック

- **Framework**: Next.js 16 (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS v4
- **Linting**: ESLint

## セットアップ

### 1. 依存関係のインストール

```bash
npm install
```

### 2. 環境変数の設定

`.env.example`をコピーして`.env.local`を作成し、必要な環境変数を設定してください。

```bash
cp .env.example .env.local
```

`.env.local`の内容:

```
NEXT_PUBLIC_APP_NAME=R2 Manager
```

## 開発

開発サーバーを起動:

```bash
npm run dev
```

ブラウザで [http://localhost:3000](http://localhost:3000) を開いてアプリケーションを確認できます。

## ビルド

本番用にビルド:

```bash
npm run build
```

ビルド後のアプリケーションを起動:

```bash
npm start
```

## リンティング

ESLintでコードをチェック:

```bash
npm run lint
```

## ディレクトリ構造

```
src/frontend/
├── app/              # Next.js App Router のページとレイアウト
│   ├── layout.tsx    # ルートレイアウト
│   ├── page.tsx      # ホームページ
│   └── globals.css   # グローバルスタイル
├── public/           # 静的ファイル
├── .env.example      # 環境変数のテンプレート
├── .env.local        # 環境変数（gitignore対象）
├── next.config.ts    # Next.js設定
├── tsconfig.json     # TypeScript設定
└── tailwind.config.ts # Tailwind CSS設定
```

## 開発ガイドライン

- TypeScriptの型安全性を活用してください
- Tailwind CSSを使用してスタイリングを行ってください
- パスエイリアス `@/*` を使用してインポートを簡潔に記述できます
