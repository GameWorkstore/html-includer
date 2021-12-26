const att = "html-include";
const attnode = "<"+att+">"

function HtmlInclude()
{
    Include(null);
}

function Include(onUpdate)
{
    var z, i, element, file, xHttp;
    z = document.getElementsByTagName("*");
    for (i = 0; i < z.length; i++)
    {
        element = z[i];
        file = element.getAttribute(att);
        if (file)
        {
            xHttp = new XMLHttpRequest();
            xHttp.responseType = 'text';
            xHttp.onreadystatechange = function ()
            {
                if (this.readyState == 4)
                {
                    if (this.status == 200)
                    {
                        if(this.responseText.startsWith(attnode)){
                            var l = this.responseText.length
                            var start = attnode.length
                            var end = l - (start*2 +1)
                            element.innerHTML = element.innerHTML + this.responseText.substring(start,end)
                        }
                        else{
                            element.innerHTML = element.innerHTML + this.responseText;
                        }
                    }
                    else if (this.status == 404)
                    {
                        element.innerHTML = "HtmlInclude: " + file + " not found.";
                    }
                    element.removeAttribute(att);
                    Include(onUpdate);
                }
            }
            xHttp.open("GET", file, true);
            xHttp.send();
            return;
        }
    }
    if(onUpdate != null) onUpdate();
};

function GetMeta(metaName)
{
    for (meta of document.getElementsByTagName('meta'))
    {
        if (meta.getAttribute('name') === metaName)
        {
            return meta.getAttribute('content');
        }
    }
    return null;
}