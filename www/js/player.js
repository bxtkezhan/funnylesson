function new_player() {
    document.querySelectorAll('.player').forEach(function(frame) {
        frame.width = frame.parentElement.clientWidth;
        frame.height = Math.round(frame.width * 9 / 16);
        // frame.style.maxWidth = '100%';
        frame.setAttribute('frameborder', '0');
        frame.setAttribute('allow',
            'accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture');
    });
}

document.addEventListener('DOMContentLoaded', function(event) {
    new_player();
    window.addEventListener("resize", function() {
        document.querySelectorAll('.player').forEach(frame => {
            frame.width = frame.parentElement.clientWidth;
            frame.height = Math.round(frame.width * 9 / 16);
        });
    });
});

function set_player(id, source, autoplay=0) {
    var frame = document.querySelector(id);
    var url = new URL(source);
    frame.setAttribute('src', `https://www.ixigua.com/iframe${url.pathname}?autoplay=${autoplay}`);
}
