<div align="center">
    <img src="./.assets/bytesentinel.png" width="250px" style="margin-left: 10px" />
</div>

<h1 align="center">
  Name of repo
</h1>

<div align="center">
    <img src="https://img.shields.io/github/downloads/bytesentinel-io/go-template/total?style=for-the-badge" />
    <img src="https://img.shields.io/github/last-commit/bytesentinel-io/go-template?color=%231BCBF2&style=for-the-badge" />
    <img src="https://img.shields.io/github/issues/bytesentinel-io/go-template?style=for-the-badge" />
</div>

<br />

`DESCRIPTION OF PROJECT`

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/Z8Z8JPE9P)

# 📝 Table of Contents <a name="table-of-contents"></a>

- [📝 Table of Contents](#table-of-contents)
- [📂 Project Structure](#project-structure)
- [🧐 About](#about)
- [🏁 Getting Started](#getting_started)
  - [💻 Prerequisites](#prerequisites)
    - [🔐 Create a secure secret key](#create-a-secure-secret-key)
  - [🚀 Installing](#installing)
- [📝 License](#license)
- [📝 Acknowledgements](#acknowledgements)

# 📂 Project Structure <a name="project-structure"></a>

```
📂 src
├── 📂 api
│   ├── 📂 app
│   │   └── 📄 api.js
│   ├── 📂 config
│   │   ├── 📄 db.js
│   │   └── 📄 node.js
│   ├── 📂 models
│   │   ├── 📄 customers.js
│   │   ├── 📄 index.js
│   │   └── 📄 users.js
│   ├── 📂 routes
│   │   ├── 📂 auth
│   │   │   ├── 📄 index.js
│   │   │   └── 📄 oauth.js
│   │   └── 📂 index.js
│   └── 📄 index.js
├── 📄 .gitignore
├── 🔒 .env
├── 📄 .env.example
├── 📄 package.json
└── 📄 README.md
```

# 🧐 About <a name="about"></a>

This project is a core API for ByteSentinel. It is a Node.js API that uses **Express.js**, **Sequelize ORM**, and **PostgreSQL**. It is a RESTful API that uses JWT for authentication and authorization. It is a boilerplate for future projects.

# 🏁 Getting Started <a name="getting_started"></a>

Follow these instructions to get a copy of the project up and running on your local machine for development and testing purposes.

## 💻 Prerequisites <a name="prerequisites"></a>

- [Node.js](https://nodejs.org/en/)
- [PostgreSQL](https://www.postgresql.org/)
- [Git](https://git-scm.com/)
- [NPM](https://www.npmjs.com/)
- [Yarn](https://yarnpkg.com/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Postman](https://www.postman.com/)
- [VSCode](https://code.visualstudio.com/)

### 🔐 Create a secure secret key <a name="create-a-secure-secret-key"></a>

```bash
openssl rand -hex 32
```

### 🚀 Installing <a name="installing"></a>

Clone the repository

```bash
git clone https://github.com/bytesentinel-io/api.git
```

Install dependencies

```bash
npm install -g yarn
yarn install
```

Create a `.env` file in the root directory and copy the contents of `.env.example` into it. Fill in the values.

```bash
cp .env.example .env
```

Start the PostgreSQL database

```bash
docker run --name bytesentinel_psql -e POSTGRES_DB=dbo.name -e POSTGRES_USER=username -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres:latest
```

Start the API

```bash
yarn dev
```

# 📝 License <a name="license"></a>

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

# 📝 Acknowledgements <a name="acknowledgements"></a>

`COMING SOON`