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
    export let data: MarketQuery

    import { createClient } from "graphql-ws"
    import { onDestroy } from "svelte"
    import { GRAPHQL_SUBSCRIPTION_URL } from "$lib/env"
    import type { MarketUpdateSubscription, MarketQuery } from "$lib/query"
    import { MarketUpdateDocumentString } from "$lib/query"
    import { browser } from "$app/env"
    import { writable } from "svelte/store"

    const products = writable(
        data.marketProducts.edges.map((edge) => edge.node)
    )

    if (browser) {
        const client = createClient({
            url: GRAPHQL_SUBSCRIPTION_URL
        })

        const unsubscribe = client.subscribe<MarketUpdateSubscription>(
            {
                query: MarketUpdateDocumentString
            },
            {
                next({ errors, data }) {
                    if (errors) {
                        console.error(errors)
                    }
                    if (data) {
                        products.update((products) => {
                            products.push(data.productOffered)
                            return products
                        })
                    }
                },
                error(error) {
                    console.error(error)
                },
                complete() {
                    console.log("Complete!")
                }
            }
        )

        onDestroy(unsubscribe)
    }
</script>

<svelte:head>
    <title>Market</title>
</svelte:head>

<div class="container mx-auto">
    <h1 class="text-3xl bold ">Market</h1>
    <pre>{JSON.stringify($products, null, 2)}</pre>
</div>

<style>
</style>
