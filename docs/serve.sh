#!/bin/bash

# Install dependencies if needed
if [ ! -d "vendor" ]; then
    echo "Installing Jekyll dependencies..."
    bundle install --path vendor/bundle
fi

# Serve the site locally
echo "Starting Jekyll server..."
echo "Visit http://localhost:4000 to view the site"
bundle exec jekyll serve --watch --livereload