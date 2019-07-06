module.exports = {
    devServer: {
        proxy: {
            '/api': {
                //target: 'http://localhost:8080',
                target: 'http://192.168.1.15:8080',
                changeOrigin: true,
                ws: true
            }
        }
    }
};
