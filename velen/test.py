#!/usr/bin/env python3
import os
import pwd
import shutil
import subprocess
import tempfile
import unittest

username = 'boruta-user'
velenPath = '/tmp/velen'

def assertRootIsOverlay(testCase, mountCmdOutput):
    roots = [x for x in [x.split() for x in mountCmdOutput.strip().split('\n')] if x[2] == '/']
    testCase.assertEqual(len(roots), 1, "more than one / in velen")
    testCase.assertEqual(roots[0][4], 'overlay', "/ in velen is not overlay")

class TestVelen(unittest.TestCase):
    def setUp(self):
        self.assertFalse(os.path.exists(velenPath), "velen in dirty state")
        self.home = pwd.getpwnam(username)[5]
        subprocess.check_call(['sudo', 'velen-prepare'])
    def test_env(self):
        env = subprocess.check_output(['sudo', 'velen-run', 'env'], universal_newlines=True).strip().split('\n')
        self.assertIn('USER={}'.format(username), env, "incorrect user in velen")
        self.assertIn('HOME={}'.format(self.home), env, "incorrect user home in velen")
    def test_uid(self):
        uid = subprocess.check_output(['sudo', 'velen-run', 'id', '-un'], universal_newlines=True).strip()
        self.assertEqual(username, uid, "incorrect uid in velen")
    def test_gid(self):
        gid = subprocess.check_output(['sudo', 'velen-run', 'id', '-gn'], universal_newlines=True).strip()
        self.assertEqual(username, gid, "incorrect gid in velen")
    def test_root_mount(self):
        mounts = subprocess.check_output(['sudo', 'velen-run', 'mount'], universal_newlines=True)
        assertRootIsOverlay(self,mounts)
    def test_shell_uid(self):
        uid = subprocess.check_output(['velen-shell', '-c', 'id -un'], universal_newlines=True).strip()
        self.assertEqual(username, uid, "incorrect uid in velen")
    def tearDown(self):
        subprocess.check_call(['sudo', 'velen-destroy'])
        self.assertFalse(os.path.exists(velenPath), "velen in dirty state after destroy")

class TestSshShell(unittest.TestCase):
    def setUp(self):
        self.assertFalse(os.path.exists(velenPath), "velen in dirty state before test")
        self.homedir = pwd.getpwnam(username)[5]
        if 'velen-shell' not in pwd.getpwnam(username)[6]:
            self.skipTest("target user's shell is not velen")
        self.keydir = tempfile.mkdtemp()
        subprocess.check_call(['ssh-keygen', '-t', 'ed25519', '-f', os.path.join(self.keydir, 'key'), '-N', ''], stderr=subprocess.DEVNULL, stdout=subprocess.DEVNULL)
        subprocess.check_call(['sudo', '-H', '-u', username, 'mkdir', '-p', '-m', '0700', os.path.join(self.homedir, '.ssh')])
        with open(os.path.join(self.keydir, 'key.pub')) as f:
            pubkey = f.read() + '\n'
        authkeys = os.path.join(self.homedir, '.ssh', 'authorized_keys')
        subprocess.run(['sudo', '-H', '-u', username, 'tee', '-a', authkeys], universal_newlines=True, input=pubkey, check=True)
        subprocess.check_call(['sudo', '-H', '-u', username, 'chmod', '0600', authkeys])
        subprocess.check_call(['sudo', 'velen-prepare'])
    @unittest.skip("TODO")
    def test_ssh_shell(self):
        proc = subprocess.run(['ssh', '-t', '-o', 'UserKnownHostsFile=/dev/null', '-o', 'StrictHostKeyChecking=no', '-i', os.path.join(self.keydir, 'key'), '{}@localhost'.format(username)], universal_newlines=True, check=True, input='mount\n', stdout=subprocess.PIPE)
        assertRootIsOverlay(self,proc.stdout)
    def test_ssh_command(self):
        proc = subprocess.check_output(['ssh', '-o', 'UserKnownHostsFile=/dev/null', '-o', 'StrictHostKeyChecking=no', '-i', os.path.join(self.keydir, 'key'), '{}@localhost'.format(username), 'mount'], universal_newlines=True)
        assertRootIsOverlay(self,proc)
    def tearDown(self):
        shutil.rmtree(self.keydir)
        subprocess.check_call(['sudo', 'velen-destroy'])
        self.assertFalse(os.path.exists(velenPath), "velen in dirty state after destroy")

class TestDoublePrepare(unittest.TestCase):
    def setUp(self):
        self.assertFalse(os.path.exists(velenPath), "velen in dirty state before test")
        subprocess.check_call(['sudo', 'velen-prepare'])
    def test_double_prepare(self):
        self.assertEqual(subprocess.call(['sudo', 'velen-prepare']), 1, "velen-prepare doesn't fail when run twice")
    def tearDown(self):
        subprocess.check_call(['sudo', 'velen-destroy'])
        self.assertFalse(os.path.exists(velenPath), "velen in dirty state after destroy")

class TestChangesStayInsideVelen(unittest.TestCase):
    def setUp(self):
        self.uid = pwd.getpwnam(username)[2]
        self.assertFalse(os.path.exists(velenPath), "velen in dirty state before test")
        self.fd, self.path = tempfile.mkstemp(dir='/home/pi') # TODO: find a directory this can be done in that isn't /tmp (which is bindmounted inside velen)
        os.close(self.fd)
        subprocess.check_call(['sudo', 'velen-prepare', self.path])
    def test_chowned(self):
        outstat = os.stat(self.path)
        self.assertEqual(outstat.st_uid, os.getuid(), "file owner changed outside velen")
        instat = os.stat(velenPath + "/overlay" + self.path)
        self.assertEqual(instat.st_uid, self.uid, "file owner not changed inside velen")
    def test_delete(self):
        subprocess.check_call(['sudo', 'velen-run', 'rm', self.path])
        self.assertTrue(os.path.exists(self.path), "deletion inside velen deleted outside file")
    def test_write(self):
        with open(self.path) as f:
            good_content = f.read()
        subprocess.check_call(['sudo', 'velen-run', 'sh', '-c', 'echo bad_content > {}'.format(self.path)])
        with open(self.path) as f:
            new_content = f.read()
        self.assertEqual(good_content, new_content, 'file contents changed outside velen')
    def tearDown(self):
        subprocess.check_call(['sudo', 'velen-destroy'])
        os.remove(self.path)

if __name__ == '__main__':
    unittest.main()
