function select_template(id) {
    var template = document.querySelector(id);
    var container = template.parentNode;
    container.removeChild(template);
    template.removeAttribute('id');
    return {tpl: template, ptr: container};
}

function extend_columns(data, columns, onclick) {
    for (var i = 0; i < columns.length; ++i) {
        var column = columns[i];
        var name = column.getAttribute('f-name');
        if (name == null) continue;
        var value = data[name];
        switch (column.tagName) {
            case 'IMG':
                column.src = `/img/${value}`;
                break;
            case 'TIME':
                const time = new Date(value * 1000).toLocaleDateString();
                column.setAttribute('datetime', time);
                column.innerText = time;
                break;
            default:
                column.innerText = value;
                break;
        }
        var param = column.getAttribute('f-onclick');
        if (param != null && onclick != null) {
            const arg = data[param];
            column.onclick = function() {
                onclick(arg);
            };
            column.style.cursor = "pointer";
        }
    }
}

function extend_items(items, template, container, onclick=null) {
    items.forEach(item => {
        var row = template.cloneNode(true);
        extend_columns(item, row.children, onclick);
        container.append(row);
    });
}

async function load_course(id) {
    const resp = await fetch(`/api/course?id=${id}`);
    return await resp.json();
}

async function load_courses(page, size, category='') {
    const resp = await fetch(`/api/courses?page=${page}&size=${size}&category=${category}`);
    return await resp.json();
}

async function load_lessons(page, size) {
    const resp = await fetch(`/api/lessons?page=${page}&size=${size}`);
    return await resp.json();
}

async function load_contents(id) {
    const resp = await fetch(`/api/contents?id=${id}`);
    return await resp.json();
}

function set_loading(id) {
    var box = document.querySelector(id);
    if (box.children.length == 0) {
        var lds = document.createElement('DIV');
        lds.classList.add('lds-default');
        box.append(lds);
        for (let i = 0; i < 12; ++i) {
            lds.append(document.createElement('DIV'));
        }
    }
}

function del_loading(id) {
    var box = document.querySelector(id);
    if (box.children.length > 0) {
        box.removeChild(box.firstChild);
    }
}

function offset_footer(tag_id, footer_id) {
    var tag = document.querySelector(tag_id);
    var footer = document.querySelector(footer_id);
    const main_top = tag.getBoundingClientRect().top;
    const footer_height = footer.getBoundingClientRect().height;
    tag.style.minHeight = window.innerHeight - main_top - footer_height + 'px';
}
