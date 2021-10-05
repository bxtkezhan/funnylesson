function innerCFW(callback) {
    if (localStorage.getItem('isInnerCFW')) {
        callback(true);
        return;
    }
    var url = 'https://graph.facebook.com/feed?callback=t';
    var xhr = new XMLHttpRequest();
    var called = false;
    xhr.open('GET', url);
    xhr.onreadystatechange = function() {
        if (xhr.readyState === 4 && xhr.status === 200) {
            called = true;
            callback(false);
        }
    };
    xhr.send();
    setTimeout(function() {
        if (!called) {
            xhr.abort();
            localStorage.setItem('isInnerCFW', 'Y');
            callback(true);
        }
    }, 1000);
}

function setPlayer(isInnerCFW) {
    document.querySelectorAll('.player').forEach(function(frame) {
        var videos = frame.getAttribute('data').split(',');
        if (!isInnerCFW) {
            frame.setAttribute('src', '//www.youtube.com/embed/' + videos[0]);
        } else {
            frame.setAttribute('src', '//player.bilibili.com/player.html?bvid=' + videos[1]);
        }
        frame.width = frame.parentElement.clientWidth;
        frame.height = Math.round(frame.width * 9 / 16);
        frame.style.maxWidth = '100%';
        frame.setAttribute('frameborder', '0');
        frame.setAttribute('allow',
            'accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture');
    });
}

document.addEventListener('DOMContentLoaded', function(event) {
    innerCFW(setPlayer);
    window.addEventListener("resize", function() {
        document.querySelectorAll('.player').forEach(frame => {
            frame.width = frame.parentElement.clientWidth;
            frame.height = Math.round(frame.width * 9 / 16);
        });
    });
});
