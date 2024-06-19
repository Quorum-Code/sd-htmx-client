// deno-lint-ignore-file
// Static file

function onload() {
    OnLoginRequest();
}

function OnLoginRequest() {
    const username = document.querySelector('.username')
    const password = document.querySelector('.password')
    const submitBtn = document.querySelector('.submit-btn')

    submitBtn.addEventListener('click', () => {
        submitBtn.disabled = true;
        fetch('/api/login', {
            method: 'post',
            headers: new Headers({'Content-Type': 'application/json'}),
            body: JSON.stringify({
                username: username.value,
                password: password.value
            })
        })
        .then(response => {
            return response.json();
        })
        .then(data => {
            if (data["status"] == 200) {
                alert(data["message"])
                location.replace("/")
                return;
            }

            alert(data["message"])
            submitBtn.disabled = false;
        })
    });
}

function getCookie(cname) {
    console.log("started getting cookie");
    let name = cname + "=";
    let decodedCookie = decodeURIComponent(document.cookie);
    let ca = decodedCookie.split(';');
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

function setCookie(cname, cvalue, exdays) {
    const d = new Date();
    d.setTime(d.getTime() + (exdays * 24 * 60 * 60 * 1000));
    let expires = "expires="+d.toUTCString();
    document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/";
}