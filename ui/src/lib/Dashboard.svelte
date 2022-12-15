<script lang="ts">
  import { onMount } from "svelte";
  import ApiClient, { ContainerResponse } from "./api";
  import ContainerCard from "./ContainerCard.svelte";

  let containers: ContainerResponse = null;
  const client = new ApiClient()

  onMount(async () => {
    containers = await client.getContainers()
  })  
</script>

<section class="px-3 sm:px-10 md:lg:xl:px-40 py-20 bg-opacity-10">
    {#if containers}
    <div class="grid grid-cols-1 sm:md:lg:grid-cols-3 xl:grid-cols-4 dark:bg-slate-800 bg-white">
        {#each containers.current as container}
            <ContainerCard {container} />
        {/each}
    </div>
    {/if}
</section>