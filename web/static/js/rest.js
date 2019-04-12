
function restDebugRequest(method, endpoint) {
    console.log(`REST API :: REQUEST :: ${method} ${endpoint}`);
}

function restDebugRespone(data, status) {
    console.log(`REST API :: RESPONSE :: [${status}]`, data);
}

// ------------------------------
// --- REQUESTS 

// GET /api/localsounds
function getLocalSounds(sortBy) {
    var url = '/api/localsounds'
    if (sortBy)
        url += '?sort=' + sortBy.toUpperCase()

    restDebugRequest('GET', url);

    return new Promise((resolve, rejects) => {
        
        $.getJSON(url, (res, s) => {
            restDebugRespone(res, s);
            if (s == 'success') {
                resolve(res.results);
            } else {
                rejects(res, s);
            }
        }).fail((e) => {
            rejects(e);
        });

    });
}