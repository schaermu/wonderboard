<script lang="ts">
  import { onMount } from "svelte";
  import ApiClient, { ContainerGroupedResponse } from "./api";
  import ContainerCard from "./ContainerCard.svelte";

  let containerGroups: ContainerGroupedResponse = null;
  const client = new ApiClient();

  onMount(async () => {
    containerGroups = await client.getContainersGrouped("project");
  });
</script>

<section class="px-3 sm:px-10 md:lg:xl:px-40 py-20 bg-opacity-10">
  {#if containerGroups}
    {#each containerGroups.groups as group}
      <div>
        <p class="text-3xl border-b-2 border-slate-700 pb-5 mx-5">{group.name}</p>
        <div
          class="grid grid-cols-1 sm:md:lg:grid-cols-3 xl:grid-cols-4 dark:bg-slate-800 bg-white"
        >
          {#each group.items as container}
            <ContainerCard {container} />
          {/each}
        </div>
      </div>
    {/each}
  {/if}
</section>
