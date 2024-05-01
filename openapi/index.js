import koa from "koa";
import koaStatic from "koa-static";

let app = new koa();

app.use(koaStatic("static/"));
app.use(koaStatic("bundled/"));

app.use((ctx, next) => {
  if (ctx.path === "/healthcheck" && ctx.method === "GET") {
    ctx.status = 200;
    return;
  }
  return next();
});

app.listen(3334, () => console.info("OpenAPI server started."));
