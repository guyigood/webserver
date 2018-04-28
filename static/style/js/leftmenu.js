function navBar(strData) {
    var data;
    if (typeof(strData) == "string") {
        var data = JSON.parse(strData); //部分用户解析出来的是字符串，转换一下
    } else {
        data = strData;
    }
    var html = "<ul class=\"layui-nav layui-nav-tree\">";
    $.each(data, function (i, item) {
        html += '<li class="layui-nav-item">';
        html += '<a href="javascript:;">';
        if (item.icon != undefined && item.icon != '') {
            if (item.icon.indexOf("icon-") != -1) {
                html += '<i class="iconfont ' + item.icon + '" data-icon="' + item.icon + '"></i>';
            } else {
                html += '<i class="layui-icon" data-icon="' + item.icon + '">' + item.icon + '</i>';
            }
        }
        html += '<cite>' + item.title + '</cite>';
        //如果有第二级菜单
        if (item.children !== undefined && item.children.length > 0) {
            html += '<span class="layui-nav-more"></span>';
            html += '</a>';
            html += '<dl class="layui-nav-child">';

            $.each(item.children, function (j, item2) {
                //如果有三级菜单
                html += "<dd>"
                if (item2.children !== undefined && item2.children.length > 0) {
                    html += '<a href="javascript:;">' ;
                    html += '<cite>' + item2.title + '</cite>';
                    html += '<span class="layui-nav-more"></span>';
                    html += '</a>';
                    html += '<dl class="layui-nav-child">';
                    $.each(item2.children, function (k, item3) {
                        html += '<dd><a href="javascript:;" data-url="' + item3.href + '">';
                        html += '<cite>' + item3.title + '</cite>';
                        html += '</a></dd>';
                    });
                    html += "</dl>";
                }
                else {
                    html += '<a href="javascript:;" data-url="' + item2.href + '">';
                    html += '<cite>' + item2.title + '</cite>';
                    if (item2.icon != undefined && item2.icon != '') {
                        if (item2.icon.indexOf("icon-") != -1) {
                            html += '<i class="iconfont ' + item2.icon + '" data-icon="' + item2.icon + '"></i>';
                        } else {
                            html += '<i class="layui-icon" data-icon="' + item2.icon + '">' + item2.icon + '</i>';
                        }
                    }

                    html += '</a>';
                }
                html+="</dd>";
            });

            html += "</dl>";
        } else {
            html += '</a>';
        }
        html += "</li>";
    });
    html += '</ul>';
    return html;
}