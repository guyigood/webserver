layui.config({
    base: "/static/style/js/"
}).use(['table','jquery', 'element','layer','laypage'], function () {
    var laypage = layui.laypage,
        element = layui.element(),
        layer=layui.layer,
        table = layui.table,
        $ = layui.jquery;
    table.render({
        elem: '#list'
        , height: 315
        , url: '/nav_getlist/' //数据接口
        , cols: [[ //表头
            {field: 'id', title: 'ID', width: 80, sort: true, fixed: 'left'}
            , {field: 'nav_name', title: '菜单名称', width: 80}
            , {field: 'parent_id', title: '上级菜单ID', width: 80, sort: true}
            , {field: 'nav_module', title: '菜单链接', width: 80}
            , {field: 'nav_image', title: '菜单图标', width: 177}
            , {field: 'is_display', title: '是否显示', width: 80, sort: true,
                templet: function (d) {
                    if (d.is_display == 1) {
                        return '<span class="layui-badge layui-bg-green">是</span>';
                    } else {
                        return '<span class="layui-badge layui-bg-gray">否</span>';
                    }
                }
            },
            {field: 'is_tel', title: '手机菜单', width: 80, sort: true,
                templet: function (d) {
                    if (d.is_tel == 1) {
                        return '<span class="layui-badge layui-bg-green">是</span>';
                    } else {
                        return '<span class="layui-badge layui-bg-gray">否</span>';
                    }
                  }
                }
                ,{field: 'order_number', title: '排序', width: 80}
        ]]
    });

    laypage.render({
        elem: 'page' //注意，这里的 test1 是 ID，不用加 # 号
        ,count: 50 //数据总数，从服务端得到
    });
});
