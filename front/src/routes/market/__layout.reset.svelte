<script lang="ts">
    import { page } from "$app/stores"
    import "./layout.css"

    const menuItems = [
        {
            path: "/",
            name: "Home"
        },
        {
            path: "/about",
            name: "About"
        },
        {
            path: "/todos",
            name: "Todos"
        },
        {
            path: "/market",
            name: "Market"
        }
    ]
</script>

<div class="bg-gradient-to-br from-slate-200 to-white">
    <header class="mx-auto h-16">
        <nav>
            <ul class="flex">
                {#each menuItems as it}
                    <li class:active={$page.url.pathname === it.path}>
                        <a sveltekit:prefetch href="/" class="">{it.name}</a>
                    </li>
                {/each}
            </ul>
        </nav>
    </header>
    <slot />
</div>
