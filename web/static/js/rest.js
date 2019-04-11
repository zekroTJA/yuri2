
function restDebugRequest(method, endpoint) {
    console.log(`REST API :: REQUEST :: ${method} ${endpoint}`);
}

function restDebugRespone(data, status) {
    console.log(`REST API :: RESPONSE :: [${status}]`, data);
}

// ------------------------------
// --- REQUESTS 

function getLocalSounds() {
    restDebugRequest('GET', '/api/localsounds');
    return new Promise((resolve, rejects) => {
        $.getJSON('/api/localsounds', (res, s) => {
            restDebugRespone(s, res);
            if (s == 'success') {
                resolve(res.results);
            } else {
                rejects(res, s);
            }
        });
    });
}