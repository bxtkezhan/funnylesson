function select_template(id) {
    var template = document.querySelector(id);
    var container = template.parentNode;
    container.removeChild(template);
    template.removeAttribute('id');
    return {tpl: template, ptr: container};
}

function extend_items(items, template, container, onclick=null) {
    items.forEach(item => {
        var row = template.cloneNode(true);
        var columns = row.children;
        for (var i = 0; i < columns.length; ++i) {
            var column = columns[i];
            var value = item[column.getAttribute('f-name')];
            switch (column.tagName) {
                case 'IMG':
                    column.src = `/img/${value}`;
                    break;
                case 'TIME':
                    const time = new Date(value * 1000).toDateString();
                    column.setAttribute('datetime', time);
                    column.innerText = time;
                    break;
                default:
                    column.innerText = value;
                    break;
            }
            var param = column.getAttribute('f-onclick');
            if (param != null && onclick != null) {
                const arg = item[param];
                column.onclick = function() {
                    onclick(arg);
                };
                column.style.cursor = "pointer";
            }
        }
        container.append(row);
    });
}

async function load_courses(page, size) {
    const resp = await fetch(`/api/courses?page=${page}&size=${size}`);
    return await resp.json();
}

async function load_lessons(page, size) {
    const resp = await fetch(`/api/lessons?page=${page}&size=${size}`);
    return await resp.json();
}