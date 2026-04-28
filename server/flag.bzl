def _flag_to_file_impl(ctx):
    out = ctx.actions.declare_file(ctx.label.name + ".txt")
    ctx.actions.write(out, ctx.build_setting_value)
    return [DefaultInfo(files = depset([out]))]

flag_to_file = rule(
    implementation = _flag_to_file_impl,
    build_setting = config.string(flag = True)
)
