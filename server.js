// server.js
const express = require('express');
const path = require('path');

const app = express();
const port = 3311;

// Serve ReDoc from CDN
app.get('/', (req, res) => {
    res.send(`
        <!DOCTYPE html>
        <html>
        <head>
            <title>API Documentation</title>
            <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.css">
        </head>
        <body>
            <redoc spec-url="/api-docs"></redoc>
            <script src="https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js"></script>
        </body>
        </html>
    `);
});

// Serve OpenAPI specification file
app.use('/api-docs', express.static(path.join(__dirname, 'docs/swagger.yaml')));

app.listen(port, () => {
    console.log(`Server is running at http://localhost:${port}`);
});
