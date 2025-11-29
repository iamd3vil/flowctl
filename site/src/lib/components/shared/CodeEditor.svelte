<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { EditorView, basicSetup } from 'codemirror';
  import { EditorState } from '@codemirror/state';
  import { python } from '@codemirror/lang-python';
  import { json } from '@codemirror/lang-json';
  import { yaml } from '@codemirror/lang-yaml';
  import { shell } from '@codemirror/legacy-modes/mode/shell';
  import { StreamLanguage } from '@codemirror/language';
  import { oneDark } from '@codemirror/theme-one-dark';
  import { keymap } from '@codemirror/view';
  import { indentWithTab } from '@codemirror/commands';

  let {
    value = $bindable(''),
    language = 'shell',
    theme = 'light',
    height = '300px',
    readonly = false,
    onchange
  }: {
    value: string;
    language?: string;
    theme?: string;
    height?: string;
    readonly?: boolean;
    onchange?: (value: string) => void;
  } = $props();

  let editorContainer: HTMLDivElement;
  let editor: EditorView | null = null;
  let currentHeight = $state(height);
  let isDragging = $state(false);
  let startY = 0;
  let startHeight = 0;

  const minHeight = 150;

  // Language options for the dropdown
  const languageOptions = [
    { value: 'shell', label: 'Shell/Bash' },
    { value: 'python', label: 'Python' },
    { value: 'yaml', label: 'YAML' },
    { value: 'json', label: 'JSON' },
    { value: 'toml', label: 'TOML' }
  ];

  function getLanguageExtension(lang: string) {
    switch (lang) {
      case 'python':
        return python();
      case 'json':
        return json();
      case 'yaml':
        return yaml();
      case 'shell':
        return StreamLanguage.define(shell);
      case 'toml':
        // TOML uses a simple text mode for now
        return [];
      default:
        return [];
    }
  }

  onMount(() => {
    if (!editorContainer) return;

    const extensions = [
      basicSetup,
      keymap.of([indentWithTab]),
      getLanguageExtension(language),
      EditorView.updateListener.of((update) => {
        if (update.docChanged && !readonly) {
          const newValue = update.state.doc.toString();
          value = newValue;
          if (onchange) {
            onchange(newValue);
          }
        }
      }),
      EditorView.editable.of(!readonly),
    ];

    if (theme === 'dark') {
      extensions.push(oneDark);
    }

    const startState = EditorState.create({
      doc: value,
      extensions
    });

    editor = new EditorView({
      state: startState,
      parent: editorContainer
    });
  });

  onDestroy(() => {
    if (editor) {
      editor.destroy();
    }
  });

  // Update editor when external value changes
  $effect(() => {
    if (editor && value !== editor.state.doc.toString()) {
      editor.dispatch({
        changes: {
          from: 0,
          to: editor.state.doc.length,
          insert: value
        }
      });
    }
  });

  // Update language when changed
  $effect(() => {
    if (editor) {
      const newExtensions = [
        basicSetup,
        keymap.of([indentWithTab]),
        getLanguageExtension(language),
        EditorView.updateListener.of((update) => {
          if (update.docChanged && !readonly) {
            const newValue = update.state.doc.toString();
            value = newValue;
            if (onchange) {
              onchange(newValue);
            }
          }
        }),
        EditorView.editable.of(!readonly),
      ];

      if (theme === 'dark') {
        newExtensions.push(oneDark);
      }

      // Recreate editor with new language
      const currentValue = editor.state.doc.toString();
      editor.destroy();

      const newState = EditorState.create({
        doc: currentValue,
        extensions: newExtensions
      });

      editor = new EditorView({
        state: newState,
        parent: editorContainer
      });
    }
  });

  function handleLanguageChange(event: Event) {
    const target = event.target as HTMLSelectElement;
    language = target.value;
  }

  function handleMouseDown(event: MouseEvent) {
    isDragging = true;
    startY = event.clientY;
    startHeight = parseInt(currentHeight);
    event.preventDefault();
  }

  function handleMouseMove(event: MouseEvent) {
    if (!isDragging) return;
    const delta = event.clientY - startY;
    const newHeight = Math.max(minHeight, startHeight + delta);
    currentHeight = `${newHeight}px`;
  }

  function handleMouseUp() {
    isDragging = false;
  }

  $effect(() => {
    if (isDragging) {
      document.addEventListener('mousemove', handleMouseMove);
      document.addEventListener('mouseup', handleMouseUp);
    } else {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
    }
  });
</script>

<div class="code-editor-wrapper">
  <!-- Language selector -->
  <div class="flex items-center justify-between mb-2">
    <div class="flex items-center gap-2">
      <label for="language-select" class="text-sm font-medium text-gray-700">
        Language:
      </label>
      <select
        id="language-select"
        value={language}
        onchange={handleLanguageChange}
        class="px-2 py-1 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
      >
        {#each languageOptions as option}
          <option value={option.value}>{option.label}</option>
        {/each}
      </select>
    </div>

    <div class="flex items-center gap-2 text-xs text-gray-500">
      <span>Ctrl+Space for suggestions</span>
    </div>
  </div>

  <!-- Editor container -->
  <div class="editor-container-wrapper">
    <div
      bind:this={editorContainer}
      class="border border-gray-300 rounded-t-md overflow-auto"
      style="height: {currentHeight}"
      onwheel={(e) => e.stopPropagation()}
    ></div>

    <!-- Resize handle -->
    <div
      class="resize-handle"
      onmousedown={handleMouseDown}
      role="separator"
      aria-orientation="horizontal"
      aria-label="Resize editor"
    >
      <div class="resize-handle-line"></div>
    </div>
  </div>
</div>

<style>
  .code-editor-wrapper {
    width: 100%;
  }

  .editor-container-wrapper {
    position: relative;
  }

  .resize-handle {
    height: 12px;
    background: #f3f4f6;
    border: 1px solid #d1d5db;
    border-top: none;
    border-radius: 0 0 0.375rem 0.375rem;
    cursor: ns-resize;
    display: flex;
    align-items: center;
    justify-content: center;
    user-select: none;
    transition: background-color 0.2s;
  }

  .resize-handle:hover {
    background: #e5e7eb;
  }

  .resize-handle:active {
    background: #d1d5db;
  }

  .resize-handle-line {
    width: 40px;
    height: 3px;
    background: #9ca3af;
    border-radius: 2px;
  }

  .resize-handle:hover .resize-handle-line {
    background: #6b7280;
  }

  /* Ensure CodeMirror editor fills the container */
  :global(.code-editor-wrapper .cm-editor) {
    height: 100%;
  }

  :global(.code-editor-wrapper .cm-scroller) {
    overflow: auto;
  }
</style>
