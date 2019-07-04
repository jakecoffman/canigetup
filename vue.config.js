module.exports = {
    devServer: {
        proxy: {
            '/api': {
                target: 'http://localhost:8080',
                // target: 'https://set.jakecoffman.com',
                changeOrigin: true,
                ws: true
            }
        }
    }
}
