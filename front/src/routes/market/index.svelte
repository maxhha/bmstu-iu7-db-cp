<script context="module">
    import { GRAPHQL_URL } from "$lib/env"
    import { GraphQLClient } from "graphql-request"
    import { getSdk } from "$lib/query"

    /** @type {import('@sveltejs/kit').Load} */
    export async function load({ fetch }) {
        const client = new GraphQLClient(GRAPHQL_URL, { fetch })
        const response = await getSdk(client).Market()

        return {
            status: response.status,
            props: {
                data: response.data
            }
        }
    }
</script>

<script lang="ts">
    export let data: any
</script>

<svelte:head>
    <title>Market</title>
</svelte:head>

<div class="market-content">
    <h1>Market</h1>
    <pre>{JSON.stringify(data, null, 2)}</pre>
</div>

<style>
</style>
