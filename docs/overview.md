## A basic example in order to get the big picture

Roughly said, `manala` allows you to embed distributed templates in your projects and ease the synchronization of your projects when the reference templates are updated.

In this usage example, we are going to implement a very basic feature, yet rich enough to measure the benefits of `manala` and to fully understand the basic concepts behind it.

## The scenario of our example

All your company's projects use [`PHP-CS-fixer`](https://github.com/FriendsOfPHP/PHP-CS-Fixer) in order to define your coding rules and apply them.

Your company would like to always apply the same coding rules on all of its projects, but maintaining the same set of rules in every project can be tedious and error-prone. In an idealistic world, the coding rules should be maintained in one place and passed on to all your projects as seamlessly as possible.

That's where manala enters the game ...

## Some wording: `project` vs `recipe` vs `repository`

In manala's vocabulary, your projects (the company's PHP projects in our example) are called ... `projects`.

In our example, our reference coding rules will be stored in a single place where they will be maintained. A `recipe` is a set of templates (the file containing your coding rules is one of these templates). All the recipes and templates you maintain are made accessible to your colleagues through a `repository`.

## First step: install manala

See [installation documentation](installation.md)

!!! Tip
    Run `manala` in a console/terminal to obtain a list of available commands.

## Create your recipe repository and your first template

!!! Tip
    manala ships with some recipes by default. Run `manala list` to get the list of available recipes.

But in this example, we are going to create our own recipe repository to better understand how manala works under the hood and enable you to develop your own recipes and templates when the need arises.

Run the following command to create your recipe repository: 

`mkdir ~/my-manala-recipe-repository`

Within this repository, we are going to create a set of templates that will host our PHP rule template:

`mkdir mkdir ~/my-manala-recipe-repository/my-php-templates`

!!! Note
    In manala's philosophy, a repository is viewed as a company-wide repository where you can store recipes and templates for various purposes and many profiles: devops, hosting, backend developers, frontend developers, etc. In fact, your projects will not embed all the company's recipes but just the subset of recipes that are useful for your project. In our example, we are going to embed only the templates under `my-php-templates` in our PHP projects.

Let's create a `.manala.yaml` file under the `my-php-templates`:

```shell
  cd ~/my-manala-recipe-repository/my-php-templates
  touch ./.manala.yaml
```

!!! Note
    the `.manala.yaml` file acts as a manifest for your recipe. It holds its description, and indicates which files or folders must be put under synchronization.

Now edit this file and put the following content:

```yaml
manala:
    description: My company's PHP recipe
    sync:
        - .manala
```

Now we are going to create the `.manala` folder where all our PHP templates will be hosted:

```shell
mkdir ./.manala
```

And finally, our PHP rule template:

```shell
touch ./.manala/php-cs-rules.php
```

And paste the following content:

```php
<?php

$header = <<<'EOF'
This file is part of the XXX project.

Copyright © My Company

@author My Company <contact@my-company.com>
EOF;

return [
    '@Symfony' => true,
    'psr0' => false,
    'phpdoc_summary' => false,
    'phpdoc_annotation_without_dot' => false,
    'phpdoc_order' => true,
    'array_syntax' => ['syntax' => 'short'],
    'ordered_imports' => true,
    'simplified_null_return' => false,
    'header_comment' => ['header' => $header],
    'yoda_style' => null,
    'native_function_invocation' => ['include' => ['@compiler_optimized']],
    'no_superfluous_phpdoc_tags' => true,
];

```

## Embed our templates in a PHP project

### Create a PHP project

For the sake of our example, we are going to create a blank PHP project, but you can of course skip this step if you already have a current PHP project that uses `PHP-CS-fixer`.

```shell
mkdir ~/my-php-project
cd ~/my-php-project
mkdir ./src
# Let's create a PHP file to give some food to PHP-CS-fixer
echo "<?php\n echo \"Coucou\";\n" > ./src/hello.php
composer init
composer require friendsofphp/php-cs-fixer
touch ./.php_cs.dist
```

Add the following content in `./php_cs.dist`:

```php
<?php

$header = <<<'EOF'
This file is part of the My-wonderful-project project.

Copyright © My company

@author My company <contact@my-company.com>
EOF;

$finder = PhpCsFixer\Finder::create()
    ->in([
        // App
        __DIR__ . '/src',
    ])
;

return PhpCsFixer\Config::create()
    ->setUsingCache(true)
    ->setRiskyAllowed(true)
    ->setFinder($finder)
    ->setRules([
        '@Symfony' => true,
        'psr0' => false,
        'phpdoc_summary' => false,
        'phpdoc_annotation_without_dot' => false,
        'phpdoc_order' => true,
        'array_syntax' => ['syntax' => 'short'],
        'ordered_imports' => true,
        'simplified_null_return' => false,
        'header_comment' => ['header' => $header],
        'yoda_style' => null,
        'native_function_invocation' => ['include' => ['@compiler_optimized']],
        'no_superfluous_phpdoc_tags' => true,
    ])
;
```

!!! Note
    For the moment, we have hard-coded our coding rules but in the next step, we will of course replace them with our shared rules.

Run `vendor/bin/php-cs-fixer fix --dry-run` to check that your PHP-CS-fix config is OK.

### Embed our PHP templates in our PHP project

Create a `.manala.yaml` at the root of your PHP project:

`touch ./.manala.yaml`

And add the following content:

```yaml
manala:
  repository: /path/to/your/home/my-manala-recipe-repository
  template: my-php-templates
```

!!! warning
     Update `/path/to/your/home/` to match your real home !!! Using `~` won't work !!!

And finally run the following command:


```shell
manala up
# More verbose:
# manala up --debug
```

This command should have created a `.manala` folder at the root of your project, including a `php-cs-rules.php` file.

### And use our shared PHP rules

Replace the content of the `.php_cs.dist` file with the following code, in order to include the rules that are defined from now on in `.manala/php-cs-rules.php`:

```php
<?php

$finder = PhpCsFixer\Finder::create()
    ->in([
        // App
        __DIR__ . '/src',
    ])
;

return PhpCsFixer\Config::create()
    ->setUsingCache(true)
    ->setRiskyAllowed(true)
    ->setFinder($finder)
    ->setRules(include('.manala/php-cs-rules.php'))
;

```

Run `vendor/bin/php-cs-fixer fix --dry-run` to test that everything is OK !

And finally run `vendor/bin/php-cs-fixer fix` to apply the coding rules.

That's done !

But, hey, just wait a minute ! What about the header of my coding rules ? My shared coding rules mention a hard-coded project name (`This file is part of the XXX project`), and I want this part to be dynamic, depending on the current project !

## Defining dynamic parts in your templates

In this chapter, we are going to define some dynamic parts in our templates and implement them in the projects embedding our recipe.

For this, we must first update our template to include some dynamic parts.

Rename the `php-cs-rules.php` to add a `tmpl` suffix:

```shell
	mv ~/my-manala-recipe-repository/my-php-templates/php-cs-rules.php ~/my-manala-recipe-repository/my-php-templates/php-cs-rules.php.tmpl
```

And update its content:

```diff
- This file is part of the XXX project.
+ This file is part of the {{ .Vars.project_name }} project.
```

!!! Tip
    Templates must be written according to [Golang template syntax](https://golang.org/pkg/text/template/), plus some sugar functions brought by [Sprig](http://masterminds.github.io/sprig/).

Now edit the `.manala.yaml` file in your PHP project to add the following line:

```
project_name: My-awesome-project
```

Run the `manala up` command and look at the changes:

```shell
manala up
cat ./.manala/php-cs-rules.php
```

!!! warning
    Don't forget to run `manala up` each time you edit the `.manala.yaml` file !!!

## Share your templates with your colleagues

As previously stated, recipes are meant to be distributed. GitHub is of course the right place to host your recipes!

So, push your recipe to GitHub and don't forget to update your manala manifest (`.manala.yaml`) in the projects that consume your recipe:

```diff
manala:
-  repository: /path/to/your/home/my-manala-recipe-repository
# Public repository
+  repository: https://github.com/my-company/manala-recipes.git
# Private repository
+  repository: git@github.com:my-company/my-manala-recipe-repository.git
```

From now on, each time you push updates to your GitHub repository, simply run `manala up` in your projects to pass on the last updates.
