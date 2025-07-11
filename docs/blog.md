---
layout: default
title: Blog
---

<section class="container" style="padding: 2rem 0;">
<h1>Blog</h1>

<div class="blog-list">
{% for post in site.posts %}
  <article class="blog-post">
    <h2><a href="{{ post.url | relative_url }}">{{ post.title }}</a></h2>
    <div class="blog-post-meta">
      {{ post.date | date: "%B %d, %Y" }} • {{ post.author | default: site.author }}
    </div>
    <p>{{ post.excerpt }}</p>
    <a href="{{ post.url | relative_url }}" class="read-more">Read more →</a>
  </article>
{% endfor %}
</div>

{% if site.posts.size == 0 %}
<div style="text-align: center; padding: 4rem 0;">
  <p style="font-size: 1.2rem; color: #666;">No blog posts yet. Check back soon!</p>
</div>
{% endif %}

</section>