function onload() {
    testIndexPost()
}

function testIndexPost() {
    const authToken = getCookie('access-token')
    const testBtn = document.querySelector('.test-btn')
    testBtn.addEventListener('click', () => {
        fetch('/api/index', {
            method: 'post',
            headers: new Headers({'Content-Type': 'application/json', 'Authorization': `Bearer ${authToken}`})
        })
        .then(response => {
            if (!response.ok) {
                console.log('response not ok...')
            }
            return response.json()
        })
        .then(data => {
            alert(JSON.stringify(data))
        });
    })
}

function getCookie(cname) {
    console.log("started getting cookie");
    const name = cname + "=";
    const decodedCookie = decodeURIComponent(document.cookie);
    const ca = decodedCookie.split(';');
    for (let i = 0; i < ca.length; i++) {
        let c = ca[i];
        while(c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            console.log("got a cookie");
            return c.substring(name.length, c.length);
        }
    }
    console.log("found no cookie");
    return "";
}