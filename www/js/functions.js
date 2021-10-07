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
        if (name != null) {
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
        }
        var param = column.getAttribute('f-onclick');
        if (param != null && onclick != null) {
            const arg = data[param];
            column.onclick = function() {
                onclick(arg);
            };
            column.style.cursor = "pointer";
        }
        var cases = column.getAttribute('f-switch');
        if (cases != null && name != null) {
            column.innerText = cases.split(';')[column.innerText];
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

async function load_user() {
    const resp = await fetch(`/api/user`);
    if (resp.status == 403) {
        location.href = '/login.html';
        return;
    }
    if (resp.ok) {
        localStorage.setItem('fl-login', true);
    }
    return await resp.json();
}

async function load_likes() {
    const resp = await fetch(`/api/likes`);
    return await resp.json();
}

async function in_likes(id) {
    const resp = await fetch(`/api/inlikes?course=${id}`);
    if (!resp.ok) return false;
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

function check_login() {
    if ((new URLSearchParams(location.search)).get('from') == 'logout') {
        localStorage.removeItem('fl-login')
    }
}

function create_menu(id, items=[]) {
    var node = document.querySelector(id);
    var link = document.createElement('A');
    link.classList.add('menu-logo');
    link.setAttribute('href', '/');
    node.append(link);
    var logo = document.createElement('IMG');
    logo.setAttribute('src', '/logo/studydou.png');
    link.append(logo);
    var head = document.createElement('A');
    head.classList.add('menu-title');
    head.setAttribute('href', '/');
    head.innerText = 'StudyDou';
    node.append(head);
    var item = document.createElement('A');
    item.classList.add('menu-item');
    if (localStorage.getItem('fl-login') != 'true') {
        item.setAttribute('href', '/login.html');
        item.innerText = '登錄|註冊';
    } else if (location.pathname == '/user.html'){
        item.setAttribute('href', '/api/logout');
        item.innerText = '登出';
    } else {
        item.setAttribute('href', '/user.html');
        item.innerText = '個人主頁';
    }
    node.append(item);
    items.forEach(item => {
        var a = document.createElement('A');
        a.classList.add('menu-item');
        a.setAttribute('href', item.href);
        a.innerText = item.text;
        node.append(a);
    });
}
