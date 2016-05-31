#!/usr/bin/ruby

# ./upload version to create a new release
# or
# ./upload -u version to update a existing release

require 'github_api'
require 'rubygems'
require 'zip'
require 'pathname'
require 'securerandom'

tmp_file = ENV['HOME'] + '/tmp/' + SecureRandom.hex(20) + '.zip'

files = Pathname.new(__dir__ + '/bin/').children

Zip::File.open(tmp_file, Zip::File::CREATE) do |zipfile|
    for file in files
        zipfile.add( Pathname.new(file).basename, file )
    end
end

if ARGV[0] == '-u'
    $tag_name = ARGV[1]
    $create_new = false
else
    $tag_name = ARGV[0]
    $create_new = true
end

github = Github.new oauth_token: ENV['GITHUB_TOKEN']

if $create_new
    github.repos.releases.create owner: 'roypur', repo: 'cloudflare-ddns', tag_name: $tag_name
end

releases = github.repos.releases.list owner: 'roypur', repo: 'cloudflare-ddns'

for value in releases
    if value['tag_name'] == tag_name
        $id = value['id']
        break
    end
end

github.repos.releases.assets.upload owner: 'roypur', repo: 'cloudflare-ddns', id: $id, filepath: tmp_file, name: $tag_name + '.zip', content_type: 'application/zip'

File.delete(tmp_file)

