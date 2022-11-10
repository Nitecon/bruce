==== BRUCE ====
Basic runtime for uniform compute environments

What does Bruce do?  Well bruce works in conjunction with serf listens to events that are propagated and runs deployments / runtime configuration within a compute environment.  Intent is to make use of concurrency for templating and serial execution of package installs to ensure proper state alignment.  Why not ansible? Because I don't want any python installs or additional deps etc.  And the ability to pull some keys / secrets out of vault and not some other crazy way of encryption... because...

Requirements:
- NO additional OS dependencies...
- Single binary (aka go binary)
- Multi platform (aka linux / mac)
- Must do package installs (at least yum & apt for now)
- Must configure templates (concurrently if possible)

