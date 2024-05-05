// deno-lint-ignore-file
// Static file
function onload() {
    SignupRequest();
}

function myFunction() {
    cookie = getCookie("acookie");

    cv = 0;
    console.log("cookie: " + cookie);
    if (cookie == "") {
        cv = 1;
    } else {
        cv = parseInt(cookie);
        cv++;
    }
    setCookie("acookie", cv, 10);

    alert("here is the cookie: " + cv);
}

function incCookie() {

}

function SignupRequest() {
    const username = document.querySelector('.username')
    const password = document.querySelector('.password')
    const submitBtn = document.querySelector('.submit-btn')

    submitBtn.addEventListener('click', () => {
        submitBtn.disabled = true;
        fetch('/api/sign-up', {
            method: 'post',
            headers: new Headers({'Content-Type': 'application/json'}),
            body: JSON.stringify({
                username: username.value,
                password: password.value
            })
        })
        .then(response => {
            if (!response.ok) {
                if (response.status == 409) {
                    alert('Username already taken...');
                }
                throw new Error('Response not ok');
            }
            return response.json();
        })
        .then(data => {
            setCookie("access-token", data["access-token"], 1)
            setCookie("refresh-token", data["refresh-token"], 60)
            location.replace("/");
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