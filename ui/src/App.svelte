<script lang="ts">
  import wandLogo from './assets/magic-wand.svg'
  import Dashboard from './lib/Dashboard.svelte';
</script>

<main>
  <div class="w-full mx-auto bg-slate-900">
    <div class="flex justify-center">
      <div class="flex-none w-9 py-5 mr-3">
        <img src={wandLogo} alt="wonderboard. your homelab dashboard, with sparkles."/>
      </div>
      <div class="py-5"><h1 class="font-serif text-4xl text-center">wonderboard</h1></div>
    </div>    
  </div>
  <Dashboard />
</main>

<style>
</style>