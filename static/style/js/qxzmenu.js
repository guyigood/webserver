function build_menu(data){
    var ulHtml="";
   for (var i=0;i<data.length;i++){
       if(data[i].children.length>0){
          ulHtml="<div class=\"layui-collapse\">" +
                      "<div class=\"layui-colla-item\">"+
                      "  <h2 class=\"layui-colla-title\">"+
                      "   <div class=\"layui-form-item\" pane=\"\">" +
                      "    <div class=\"layui-input-block\">\n" +
                      "      <input type=\"checkbox\" name=\"qxz[]\" value=\""+data[i].id+"\" lay-skin=\"primary\" title=\""+data[i].title+"\" checked=\"\">\n" +
                      "    </div>\n" +
                      "  </div>" +
                      "</h2>" +
                        " <div class=\"layui-colla-content layui-show\">";
           ulHtml+=build_menu(data);
           ulHtml+="</div></div>";
           return ulHtml;
       }else{
           ulHtml= "<input type=\"checkbox\" name=\"qxz[]\"  value=\""+data[i].id+"\" lay-skin=\"primary\" title=\""+data[i].title+"\" checked=\"\">\n";
       }
   }
   return ulHtml;
}
