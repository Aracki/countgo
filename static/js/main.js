document.addEventListener("DOMContentLoaded", function () {
    // ping golang unique visitor counter
    getRequest(location.origin + "/count", function (request) {
        var response = request.currentTarget.response || request.target.responseText;
        var counter = document.getElementById("counter_text");

        if (counter != null) {
            document.getElementById("counter_text")
                .innerHTML = "[" + response + " unique visitors]";
        }
    })
});

function getRequest(url, success) {
    var xhr = new XMLHttpRequest();
    xhr.open('GET', url);
    xhr.onload = success;
    xhr.send();
    return xhr;
}


window.onload = function() {
    const textElement = document.getElementById("greeting");
    const oldText = textElement.textContent;
    const newText = oldText.replace('Hello', getCurrentTime);
    textElement.textContent = newText; 
};

var getCurrentTime = function() {
    var date = new Date();
    var hours =  date.getHours();
    var minutes =  date.getMinutes();
    var current = hours + (minutes * .01);
    if (current >= 5 && current < 9){
        return 'Good morning'
    } else {
        return 'Hello';
    }
};

const changeBgButton = document.getElementById('toggle-img');
const body = document.body;

changeBgButton.addEventListener('click', function() {
    body.style.background = "url('../static/img/dunes/blue-silver.webp') no-repeat center center fixed";
});